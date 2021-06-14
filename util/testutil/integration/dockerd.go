package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/moby/buildkit/util/testutil/dockerd"
	"golang.org/x/sync/errgroup"
)

type testLogger struct {
	Name string
	Logs map[string]*bytes.Buffer
}

func (l testLogger) Logf(format string, v ...interface{}) {
	if buf, ok := l.Logs[l.Name]; !ok || buf == nil {
		l.Logs[l.Name] = &bytes.Buffer{}
	}
	fmt.Fprintf(l.Logs[l.Name], format, v...)
}

func withTestLogger(name string, logs map[string]*bytes.Buffer) dockerd.Option {
	return func(d *dockerd.Daemon) {
		d.Log = testLogger{Name: name, Logs: logs}
	}
}

// InitDockerdWorker registers a dockerd worker with the global registry.
func InitDockerdWorker() {
	Register(&dockerdWorker{})
}

type dockerdWorker struct{}

func (c dockerdWorker) Name() string {
	const dockerdBinary = "dockerd"
	return dockerdBinary
}

func (c dockerdWorker) New(ctx context.Context, cfg *BackendConfig) (b Backend, cl func() error, err error) {
	if err := requireRoot(); err != nil {
		return nil, nil, err
	}

	deferF := &multiCloser{}
	cl = deferF.F()

	defer func() {
		if err != nil {
			deferF.F()()
			cl = nil
		}
	}()

	var proxyGroup errgroup.Group
	deferF.append(proxyGroup.Wait)

	workDir, err := os.MkdirTemp("", "integration")
	if err != nil {
		return nil, nil, err
	}

	d, err := dockerd.NewDaemon(
		workDir,
		withTestLogger("creatingDaemon", cfg.Logs),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("new daemon error: %q, %s", err, formatLogs(cfg.Logs))
	}

	err = d.StartWithError()
	if err != nil {
		return nil, nil, err
	}
	deferF.append(d.StopWithError)

	logs := map[string]*bytes.Buffer{}
	if err := waitUnix(d.Sock(), 5*time.Second); err != nil {
		return nil, nil, fmt.Errorf("dockerd did not start up: %q, %s", err, formatLogs(logs))
	}

	ctx, cancel := context.WithCancel(context.Background())
	deferF.append(func() error { cancel(); return nil })

	dockerAPI, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(d.Sock()))
	if err != nil {
		return nil, nil, err
	}
	deferF.append(dockerAPI.Close)

	err = waitForAPI(ctx, dockerAPI, 5*time.Second)
	if err != nil {
		return nil, nil, err
	}

	// Create a file descriptor to be used as a Unix domain socket.
	// Remove it immediately (the name will still be valid for the socket) so that
	// we don't leave files all over the users tmp tree.
	f, err := os.CreateTemp("", "buildkit-integration")
	if err != nil {
		return
	}
	localPath := f.Name()
	f.Close()
	os.Remove(localPath)

	listener, err := net.Listen("unix", localPath)
	if err != nil {
		return nil, nil, err
	}
	deferF.append(listener.Close)

	proxyGroup.Go(func() error {
		for {
			tmpConn, err := listener.Accept()
			if err != nil {
				// Ignore the error from accept which is always a system error.
				return nil
			}
			conn, err := dockerAPI.DialHijack(ctx, "/grpc", "h2c", nil)
			if err != nil {
				return err
			}

			proxyGroup.Go(func() error {
				_, err := io.Copy(conn, tmpConn)
				if err != nil {
					return err
				}
				return tmpConn.Close()
			})
			proxyGroup.Go(func() error {
				_, err := io.Copy(tmpConn, conn)
				if err != nil {
					return err
				}
				return conn.Close()
			})
		}
	})

	return backend{
		address:   "unix://" + listener.Addr().String(),
		rootless:  false,
		isDockerd: true,
	}, cl, nil
}

func waitForAPI(ctx context.Context, apiClient *client.Client, d time.Duration) error {
	step := 50 * time.Millisecond
	i := 0
	for {
		if _, err := apiClient.Ping(ctx); err == nil {
			break
		}
		i++
		if time.Duration(i)*step > d {
			return fmt.Errorf("failed to connect to /_ping endpoint")
		}
		time.Sleep(step)
	}
	return nil
}

func SkipIfDockerd(t *testing.T, sb Sandbox, reason ...string) {
	t.Helper()
	sbx, ok := sb.(*sandbox)
	if !ok {
		t.Fatalf("invalid sandbox type %T", sb)
	}
	b, ok := sbx.Backend.(backend)
	if !ok {
		t.Fatalf("invalid backend type %T", b)
	}
	if b.isDockerd {
		t.Skipf("dockerd worker can not currently run this test due to missing features (%s)", strings.Join(reason, ", "))
	}
}
