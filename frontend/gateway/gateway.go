package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/containerd/containerd/v2/core/mount"
	"github.com/containerd/containerd/v2/defaults"
	"github.com/distribution/reference"
	apitypes "github.com/moby/buildkit/api/types"
	"github.com/moby/buildkit/cache"
	cacheutil "github.com/moby/buildkit/cache/util"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/sourceresolver"
	"github.com/moby/buildkit/executor"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend"
	"github.com/moby/buildkit/frontend/dockerui"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/frontend/gateway/container"
	"github.com/moby/buildkit/frontend/gateway/forwarder"
	pb "github.com/moby/buildkit/frontend/gateway/pb"
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/snapshot"
	"github.com/moby/buildkit/solver"
	"github.com/moby/buildkit/solver/errdefs"
	llberrdefs "github.com/moby/buildkit/solver/llbsolver/errdefs"
	opspb "github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/apicaps"
	"github.com/moby/buildkit/util/bklog"
	"github.com/moby/buildkit/util/grpcerrors"
	"github.com/moby/buildkit/util/stack"
	"github.com/moby/buildkit/util/tracing"
	"github.com/moby/buildkit/worker"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
	"github.com/moby/sys/signal"
	"github.com/moby/sys/user"
	digest "github.com/opencontainers/go-digest"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/sync/errgroup"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

const (
	keyDevel = "gateway-devel"
)

func NewGatewayFrontend(workers worker.Infos, allowedRepositories []string) (frontend.Frontend, error) {
	var parsedAllowedRepositories []string

	for _, allowedRepository := range allowedRepositories {
		sourceRef, err := reference.ParseNormalizedNamed(allowedRepository)
		if err != nil {
			return nil, err
		}
		parsedAllowedRepositories = append(parsedAllowedRepositories, reference.TrimNamed(sourceRef).Name())
	}

	return &gatewayFrontend{
		workers:             workers,
		allowedRepositories: parsedAllowedRepositories,
	}, nil
}

type gatewayFrontend struct {
	workers             worker.Infos
	allowedRepositories []string
}

func filterPrefix(opts map[string]string, pfx string) map[string]string {
	m := map[string]string{}
	for k, v := range opts {
		if after, ok := strings.CutPrefix(k, pfx); ok {
			m[after] = v
		}
	}
	return m
}

func (gf *gatewayFrontend) checkSourceIsAllowed(source string) error {
	// Returns nil if the source is allowed.
	// Returns an error if the source is not allowed.
	if len(gf.allowedRepositories) == 0 {
		// No source restrictions in place
		return nil
	}

	sourceRef, err := reference.ParseNormalizedNamed(source)
	if err != nil {
		return err
	}

	taglessSource := reference.TrimNamed(sourceRef).Name()

	if slices.Contains(gf.allowedRepositories, taglessSource) {
		// Allowed
		return nil
	}
	return errors.Errorf("'%s' is not an allowed gateway source", source)
}

func (gf *gatewayFrontend) Solve(ctx context.Context, llbBridge frontend.FrontendLLBBridge, exec executor.Executor, opts map[string]string, inputs map[string]*opspb.Definition, sid string, sm *session.Manager) (*frontend.Result, error) {
	source, ok := opts[frontend.KeySource]
	if !ok {
		return nil, errors.Errorf("no source specified for gateway")
	}

	_, isDevel := opts[keyDevel]
	var img dockerspec.DockerOCIImage
	var mfstDigest digest.Digest
	var rootFS cache.MutableRef
	var readonly bool // TODO: try to switch to read-only by default.

	var frontendDef *opspb.Definition

	err := gf.checkSourceIsAllowed(source)
	if err != nil {
		return nil, err
	}

	if isDevel {
		devRes, err := llbBridge.Solve(ctx,
			frontend.SolveRequest{
				Frontend:       source,
				FrontendOpt:    filterPrefix(opts, "gateway-"),
				FrontendInputs: inputs,
			}, "gateway:"+sid)
		if err != nil {
			return nil, err
		}
		defer func() {
			ctx := context.WithoutCancel(ctx)
			devRes.EachRef(func(ref solver.ResultProxy) error {
				return ref.Release(ctx)
			})
		}()
		if devRes.Ref == nil {
			return nil, errors.Errorf("development gateway didn't return default result")
		}
		frontendDef = devRes.Ref.Definition()
		res, err := devRes.Ref.Result(ctx)
		if err != nil {
			return nil, err
		}

		workerRef, ok := res.Sys().(*worker.WorkerRef)
		if !ok {
			return nil, errors.Errorf("invalid ref: %T", res.Sys())
		}

		rootFS, err = workerRef.Worker.CacheManager().New(ctx, workerRef.ImmutableRef, session.NewGroup(sid))
		if err != nil {
			return nil, err
		}
		defer rootFS.Release(context.TODO())
		config, ok := devRes.Metadata[exptypes.ExporterImageConfigKey]
		if ok {
			if err := json.Unmarshal(config, &img); err != nil {
				return nil, err
			}
		}
	} else {
		c, err := forwarder.LLBBridgeToGatewayClient(ctx, llbBridge, exec, opts, inputs, gf.workers, sid, sm)
		if err != nil {
			return nil, err
		}
		dc, err := dockerui.NewClient(c)
		if err != nil {
			return nil, err
		}
		nc, err := dc.NamedContext(source, dockerui.ContextOpt{
			CaptureDigest: &mfstDigest,
		})
		if err != nil {
			return nil, err
		}
		var st *llb.State
		if nc != nil {
			var dockerImage *dockerspec.DockerOCIImage
			st, dockerImage, err = nc.Load(ctx)
			if err != nil {
				return nil, err
			}
			if dockerImage != nil {
				img = *dockerImage
			}
		}
		if st == nil {
			sourceRef, err := reference.ParseNormalizedNamed(source)
			if err != nil {
				return nil, err
			}

			imr := sourceresolver.NewImageMetaResolver(llbBridge)
			ref, dgst, config, err := imr.ResolveImageConfig(ctx, reference.TagNameOnly(sourceRef).String(), sourceresolver.Opt{})
			if err != nil {
				return nil, err
			}

			sourceRef, err = reference.ParseNormalizedNamed(ref)
			if err != nil {
				return nil, err
			}

			mfstDigest = dgst

			if err := json.Unmarshal(config, &img); err != nil {
				return nil, err
			}

			if dgst != "" {
				sourceRef, err = reference.WithDigest(sourceRef, dgst)
				if err != nil {
					return nil, err
				}
			}

			src := llb.Image(sourceRef.String(), &markTypeFrontend{})
			st = &src
		}

		def, err := st.Marshal(ctx)
		if err != nil {
			return nil, err
		}

		res, err := llbBridge.Solve(ctx, frontend.SolveRequest{
			Definition: def.ToPB(),
		}, sid)
		if err != nil {
			return nil, err
		}
		defer func() {
			ctx := context.WithoutCancel(ctx)
			res.EachRef(func(ref solver.ResultProxy) error {
				return ref.Release(ctx)
			})
		}()
		if res.Ref == nil {
			return nil, errors.Errorf("gateway source didn't return default result")
		}
		frontendDef = res.Ref.Definition()
		r, err := res.Ref.Result(ctx)
		if err != nil {
			return nil, err
		}
		workerRef, ok := r.Sys().(*worker.WorkerRef)
		if !ok {
			return nil, errors.Errorf("invalid ref: %T", r.Sys())
		}
		rootFS, err = workerRef.Worker.CacheManager().New(ctx, workerRef.ImmutableRef, session.NewGroup(sid))
		if err != nil {
			return nil, err
		}
		defer rootFS.Release(context.TODO())
	}

	args := []string{"/run"}
	env := []string{}
	cwd := "/"
	if img.Config.Env != nil {
		env = img.Config.Env
	}
	if img.Config.Entrypoint != nil {
		args = img.Config.Entrypoint
	}
	if img.Config.WorkingDir != "" {
		cwd = img.Config.WorkingDir
	}
	i := 0
	for k, v := range opts {
		env = append(env, fmt.Sprintf("BUILDKIT_FRONTEND_OPT_%d", i)+"="+k+"="+v)
		i++
	}

	env = append(env, "BUILDKIT_SESSION_ID="+sid)

	dt, err := json.Marshal(gf.workers.WorkerInfos())
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal workers array")
	}
	env = append(env, "BUILDKIT_WORKERS="+string(dt))

	env = append(env, "BUILDKIT_EXPORTEDPRODUCT="+apicaps.ExportedProduct)

	meta := executor.Meta{
		Env:                       env,
		Args:                      args,
		Cwd:                       cwd,
		ReadonlyRootFS:            readonly,
		RemoveMountStubsRecursive: true,
	}

	if v, ok := img.Config.Labels["moby.buildkit.frontend.network.none"]; ok {
		if ok, _ := strconv.ParseBool(v); ok {
			meta.NetMode = opspb.NetMode_NONE
		}
	}

	curCaps := getCaps(img.Config.Labels["moby.buildkit.frontend.caps"])
	addCapsForKnownFrontends(curCaps, mfstDigest)
	reqCaps := getCaps(opts["frontend.caps"])
	if len(inputs) > 0 {
		reqCaps["moby.buildkit.frontend.inputs"] = struct{}{}
	}

	for c := range reqCaps {
		if _, ok := curCaps[c]; !ok {
			return nil, stack.Enable(grpcerrors.WrapCode(errdefs.NewUnsupportedFrontendCapError(c), codes.Unimplemented))
		}
	}

	lbf, ctx := serveLLBBridgeForwarder(ctx, llbBridge, exec, gf.workers, inputs, sid, sm)
	defer lbf.conn.Close()
	defer lbf.Discard()

	mdmnt, release, err := metadataMount(frontendDef)
	if err != nil {
		return nil, err
	}
	if release != nil {
		defer release()
	}
	var mnts []executor.Mount
	if mdmnt != nil {
		mnts = append(mnts, *mdmnt)
	}

	_, err = exec.Run(ctx, "", container.MountWithSession(rootFS, session.NewGroup(sid)), mnts, executor.ProcessInfo{Meta: meta, Stdin: lbf.Stdin, Stdout: lbf.Stdout, Stderr: os.Stderr}, nil)
	if err != nil {
		if errdefs.IsCanceled(ctx, err) && lbf.isErrServerClosed {
			err = errors.Errorf("frontend grpc server closed unexpectedly")
		}
		// An existing error (set via Return rpc) takes
		// precedence over this error, which in turn takes
		// precedence over a success reported via Return.
		lbf.mu.Lock()
		if lbf.err == nil {
			lbf.result = nil
			lbf.err = err
		}
		lbf.mu.Unlock()
	}

	return lbf.Result()
}

func metadataMount(def *opspb.Definition) (*executor.Mount, func(), error) {
	dt, err := def.MarshalVT()
	if err != nil {
		return nil, nil, err
	}
	dir, err := os.MkdirTemp("", "buildkit-metadata")
	if err != nil {
		return nil, nil, err
	}

	if err := os.WriteFile(filepath.Join(dir, "frontend.bin"), dt, 0400); err != nil {
		return nil, nil, err
	}

	return &executor.Mount{
			Src:      &bind{dir},
			Dest:     "/run/config/buildkit/metadata",
			Readonly: true,
		}, func() {
			os.RemoveAll(dir)
		}, nil
}

type bind struct {
	dir string
}

func (b *bind) Mount(ctx context.Context, readonly bool) (snapshot.Mountable, error) {
	return &bindMount{b.dir, readonly}, nil
}

type bindMount struct {
	dir      string
	readonly bool
}

func (b *bindMount) Mount() ([]mount.Mount, func() error, error) {
	return []mount.Mount{{
		Type:    "bind",
		Source:  b.dir,
		Options: []string{"bind", "ro", "nosuid", "nodev", "noexec"},
	}}, func() error { return nil }, nil
}

func (b *bindMount) IdentityMapping() *user.IdentityMapping {
	return nil
}

func (lbf *llbBridgeForwarder) Discard() {
	lbf.mu.Lock()
	defer lbf.mu.Unlock()

	for ctr := range lbf.ctrs {
		lbf.ReleaseContainer(context.TODO(), &pb.ReleaseContainerRequest{
			ContainerID: ctr,
		})
	}

	for id, workerRef := range lbf.workerRefByID {
		workerRef.Release(context.TODO())
		delete(lbf.workerRefByID, id)
	}
	if lbf.err != nil && lbf.result != nil {
		lbf.result.EachRef(func(r solver.ResultProxy) error {
			r.Release(context.TODO())
			return nil
		})
	}
	for _, r := range lbf.refs {
		r.Release(context.TODO())
	}
	lbf.refs = map[string]solver.ResultProxy{}
}

func (lbf *llbBridgeForwarder) Done() <-chan struct{} {
	return lbf.doneCh
}

func (lbf *llbBridgeForwarder) setResult(r *frontend.Result, err error) (*pb.ReturnResponse, error) {
	lbf.mu.Lock()
	defer lbf.mu.Unlock()

	if (r == nil) == (err == nil) {
		return nil, errors.New("gateway return must be either result or err")
	}

	if lbf.result != nil || lbf.err != nil {
		return nil, errors.New("gateway result is already set")
	}

	lbf.result = r
	lbf.err = err
	close(lbf.doneCh)
	return &pb.ReturnResponse{}, nil
}

func (lbf *llbBridgeForwarder) Result() (*frontend.Result, error) {
	lbf.mu.Lock()
	defer lbf.mu.Unlock()

	if lbf.result == nil && lbf.err == nil {
		return nil, errors.New("no result for incomplete build")
	}

	if lbf.err != nil {
		return nil, lbf.err
	}

	return lbf.result, nil
}

func NewBridgeForwarder(ctx context.Context, llbBridge frontend.FrontendLLBBridge, exec executor.Executor, workers worker.Infos, inputs map[string]*opspb.Definition, sid string, sm *session.Manager) LLBBridgeForwarder {
	return newBridgeForwarder(ctx, llbBridge, exec, workers, inputs, sid, sm)
}

func newBridgeForwarder(ctx context.Context, llbBridge frontend.FrontendLLBBridge, exec executor.Executor, workers worker.Infos, inputs map[string]*opspb.Definition, sid string, sm *session.Manager) *llbBridgeForwarder {
	lbf := &llbBridgeForwarder{
		callCtx:       ctx,
		llbBridge:     llbBridge,
		refs:          map[string]solver.ResultProxy{},
		workerRefByID: map[string]*worker.WorkerRef{},
		doneCh:        make(chan struct{}),
		pipe:          newPipe(),
		workers:       workers,
		inputs:        inputs,
		sid:           sid,
		sm:            sm,
		ctrs:          map[string]gwclient.Container{},
		executor:      exec,
	}
	return lbf
}

func serveLLBBridgeForwarder(ctx context.Context, llbBridge frontend.FrontendLLBBridge, exec executor.Executor, workers worker.Infos, inputs map[string]*opspb.Definition, sid string, sm *session.Manager) (*llbBridgeForwarder, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	lbf := newBridgeForwarder(ctx, llbBridge, exec, workers, inputs, sid, sm)
	serverOpt := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcerrors.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpcerrors.StreamServerInterceptor),
		grpc.MaxRecvMsgSize(defaults.DefaultMaxRecvMsgSize),
		grpc.MaxSendMsgSize(defaults.DefaultMaxSendMsgSize),
	}
	server := grpc.NewServer(serverOpt...)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	pb.RegisterLLBBridgeServer(server, lbf)

	go func() {
		serve(ctx, server, lbf.conn)
		select {
		case <-ctx.Done():
		default:
			lbf.isErrServerClosed = true
		}
		cancel(errors.WithStack(context.Canceled))
	}()

	return lbf, ctx
}

type pipe struct {
	Stdin  io.ReadCloser
	Stdout io.WriteCloser
	conn   net.Conn
}

func newPipe() *pipe {
	pr1, pw1, _ := os.Pipe()
	pr2, pw2, _ := os.Pipe()
	return &pipe{
		Stdin:  pr1,
		Stdout: pw2,
		conn: &conn{
			Reader: pr2,
			Writer: pw1,
			Closer: pw2,
		},
	}
}

type conn struct {
	io.Reader
	io.Writer
	io.Closer
}

func (s *conn) LocalAddr() net.Addr {
	return dummyAddr{}
}

func (s *conn) RemoteAddr() net.Addr {
	return dummyAddr{}
}

func (s *conn) SetDeadline(t time.Time) error {
	return nil
}

func (s *conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (s *conn) SetWriteDeadline(t time.Time) error {
	return nil
}

type dummyAddr struct{}

func (d dummyAddr) Network() string {
	return "pipe"
}

func (d dummyAddr) String() string {
	return "localhost"
}

type LLBBridgeForwarder interface {
	pb.LLBBridgeServer
	Done() <-chan struct{}
	Result() (*frontend.Result, error)
	Discard()
}

type llbBridgeForwarder struct {
	mu            sync.Mutex
	callCtx       context.Context
	llbBridge     frontend.FrontendLLBBridge
	refs          map[string]solver.ResultProxy
	workerRefByID map[string]*worker.WorkerRef
	// lastRef      solver.CachedResult
	// lastRefs     map[string]solver.CachedResult
	// err          error
	doneCh            chan struct{} // closed when result or err become valid through a call to a Return
	result            *frontend.Result
	err               error
	workers           worker.Infos
	inputs            map[string]*opspb.Definition
	isErrServerClosed bool
	sid               string
	sm                *session.Manager
	executor          executor.Executor
	*pipe
	ctrs   map[string]gwclient.Container
	ctrsMu sync.Mutex
}

func (lbf *llbBridgeForwarder) ResolveSourceMeta(ctx context.Context, req *pb.ResolveSourceMetaRequest) (*pb.ResolveSourceMetaResponse, error) {
	if req.Source == nil {
		return nil, status.Error(codes.InvalidArgument, "source is required")
	}

	ctx = tracing.ContextWithSpanFromContext(ctx, lbf.callCtx)
	var platform *ocispecs.Platform
	if p := req.Platform; p != nil {
		platform = &ocispecs.Platform{
			OS:           p.OS,
			Architecture: p.Architecture,
			Variant:      p.Variant,
			OSVersion:    p.OSVersion,
			OSFeatures:   p.OSFeatures,
		}
	}
	resolveopt := sourceresolver.Opt{
		LogName:        req.LogName,
		SourcePolicies: req.SourcePolicies,
		Platform:       platform,
	}
	resolveopt.ImageOpt = &sourceresolver.ResolveImageOpt{
		ResolveMode: req.ResolveMode,
	}
	resp, err := lbf.llbBridge.ResolveSourceMetadata(ctx, req.Source, resolveopt)
	if err != nil {
		return nil, err
	}

	r := &pb.ResolveSourceMetaResponse{
		Source: resp.Op,
	}

	if resp.Image != nil {
		r.Image = &pb.ResolveSourceImageResponse{
			Digest: string(resp.Image.Digest),
			Config: resp.Image.Config,
		}
	}
	return r, nil
}

func (lbf *llbBridgeForwarder) ResolveImageConfig(ctx context.Context, req *pb.ResolveImageConfigRequest) (*pb.ResolveImageConfigResponse, error) {
	ctx = tracing.ContextWithSpanFromContext(ctx, lbf.callCtx)
	var platform *ocispecs.Platform
	if p := req.Platform; p != nil {
		platform = &ocispecs.Platform{
			OS:           p.OS,
			Architecture: p.Architecture,
			Variant:      p.Variant,
			OSVersion:    p.OSVersion,
			OSFeatures:   p.OSFeatures,
		}
	}
	imr := sourceresolver.NewImageMetaResolver(lbf.llbBridge)
	resolveopt := sourceresolver.Opt{
		LogName:        req.LogName,
		SourcePolicies: req.SourcePolicies,
		Platform:       platform,
	}
	if sourceresolver.ResolverType(req.ResolverType) == sourceresolver.ResolverTypeRegistry {
		resolveopt.ImageOpt = &sourceresolver.ResolveImageOpt{
			ResolveMode: req.ResolveMode,
		}
	} else if sourceresolver.ResolverType(req.ResolverType) == sourceresolver.ResolverTypeOCILayout {
		resolveopt.OCILayoutOpt = &sourceresolver.ResolveOCILayoutOpt{
			Store: sourceresolver.ResolveImageConfigOptStore{
				SessionID: req.SessionID,
				StoreID:   req.StoreID,
			},
		}
	}

	ref, dgst, dt, err := imr.ResolveImageConfig(ctx, req.Ref, resolveopt)
	if err != nil {
		return nil, err
	}
	return &pb.ResolveImageConfigResponse{
		Ref:    ref,
		Digest: string(dgst),
		Config: dt,
	}, nil
}

func (lbf *llbBridgeForwarder) wrapSolveError(solveErr error) error {
	var (
		ee       *llberrdefs.ExecError
		fae      *llberrdefs.FileActionError
		sce      *solver.SlowCacheError
		inputIDs []string
		mountIDs []string
		subject  errdefs.IsSolve_Subject
	)
	if errors.As(solveErr, &ee) {
		var err error
		inputIDs, err = lbf.registerResultIDs(ee.Inputs...)
		if err != nil {
			return err
		}
		mountIDs, err = lbf.registerResultIDs(ee.Mounts...)
		if err != nil {
			return err
		}
	}
	if errors.As(solveErr, &fae) {
		subject = fae.ToSubject()
	}
	if errors.As(solveErr, &sce) {
		var err error
		inputIDs, err = lbf.registerResultIDs(sce.Result)
		if err != nil {
			return err
		}
		subject = sce.ToSubject()
	}
	return errdefs.WithSolveError(solveErr, subject, inputIDs, mountIDs)
}

func (lbf *llbBridgeForwarder) registerResultIDs(results ...solver.Result) (ids []string, err error) {
	lbf.mu.Lock()
	defer lbf.mu.Unlock()

	ids = make([]string, len(results))
	for i, res := range results {
		if res == nil {
			continue
		}
		workerRef, ok := res.Sys().(*worker.WorkerRef)
		if !ok {
			return ids, errors.Errorf("unexpected type for result, got %T", res.Sys())
		}
		ids[i] = workerRef.ID()
		lbf.workerRefByID[workerRef.ID()] = workerRef
	}
	return ids, nil
}

func (lbf *llbBridgeForwarder) Solve(ctx context.Context, req *pb.SolveRequest) (*pb.SolveResponse, error) {
	var cacheImports []frontend.CacheOptionsEntry
	for _, e := range req.CacheImports {
		if e == nil {
			return nil, errors.Errorf("invalid nil cache import")
		}
		cacheImports = append(cacheImports, frontend.CacheOptionsEntry{
			Type:  e.Type,
			Attrs: e.Attrs,
		})
	}

	for _, p := range req.SourcePolicies {
		if p == nil {
			return nil, errors.Errorf("invalid nil source policy")
		}
	}

	ctx = tracing.ContextWithSpanFromContext(ctx, lbf.callCtx)
	res, err := lbf.llbBridge.Solve(ctx, frontend.SolveRequest{
		Evaluate:       req.Evaluate,
		Definition:     req.Definition,
		Frontend:       req.Frontend,
		FrontendOpt:    req.FrontendOpt,
		FrontendInputs: req.FrontendInputs,
		CacheImports:   cacheImports,
		SourcePolicies: req.SourcePolicies,
	}, lbf.sid)
	if err != nil {
		return nil, lbf.wrapSolveError(err)
	}

	if len(res.Refs) > 0 && !req.AllowResultReturn {
		// this should never happen because old client shouldn't make a map request
		return nil, errors.Errorf("solve did not return default result")
	}

	pbRes := &pb.Result{
		Metadata: res.Metadata,
	}
	var defaultID string

	lbf.mu.Lock()

	if res.Refs != nil {
		ids := make(map[string]string, len(res.Refs))
		defs := make(map[string]*opspb.Definition, len(res.Refs))
		for k, ref := range res.Refs {
			var id string
			var def *opspb.Definition
			if ref != nil {
				id = identity.NewID()
				def = ref.Definition()
				lbf.refs[id] = ref
			}
			ids[k] = id
			defs[k] = def
		}

		if req.AllowResultArrayRef {
			refMap := make(map[string]*pb.Ref, len(res.Refs))
			for k, id := range ids {
				refMap[k] = &pb.Ref{Id: id, Def: defs[k]}
			}
			pbRes.Result = &pb.Result_Refs{Refs: &pb.RefMap{Refs: refMap}}
		} else {
			pbRes.Result = &pb.Result_RefsDeprecated{RefsDeprecated: &pb.RefMapDeprecated{Refs: ids}}
		}
	} else {
		ref := res.Ref
		var id string
		var def *opspb.Definition
		if ref != nil {
			id = identity.NewID()
			def = ref.Definition()
			lbf.refs[id] = ref
		}
		defaultID = id

		if req.AllowResultArrayRef {
			pbRes.Result = &pb.Result_Ref{Ref: &pb.Ref{Id: id, Def: def}}
		} else {
			pbRes.Result = &pb.Result_RefDeprecated{RefDeprecated: id}
		}
	}

	if res.Attestations != nil {
		pbRes.Attestations = map[string]*pb.Attestations{}
		for k, atts := range res.Attestations {
			for _, att := range atts {
				pbAtt, err := gwclient.AttestationToPB(&att)
				if err != nil {
					return nil, err
				}

				if att.Ref != nil {
					id := identity.NewID()
					def := att.Ref.Definition()
					lbf.refs[id] = att.Ref
					pbAtt.Ref = &pb.Ref{Id: id, Def: def}
				}

				if pbRes.Attestations[k] == nil {
					pbRes.Attestations[k] = &pb.Attestations{}
				}
				pbRes.Attestations[k].Attestation = append(pbRes.Attestations[k].Attestation, pbAtt)
			}
		}
	}

	lbf.mu.Unlock()

	// compatibility mode for older clients
	if req.Final {
		exp := map[string][]byte{}
		if err := json.Unmarshal(req.ExporterAttr, &exp); err != nil {
			return nil, err
		}
		maps.Copy(exp, res.Metadata)

		lbf.mu.Lock()
		lbf.result = &frontend.Result{
			Ref:      lbf.refs[defaultID],
			Metadata: exp,
		}
		lbf.mu.Unlock()
	}

	resp := &pb.SolveResponse{
		Result: pbRes,
	}

	if !req.AllowResultReturn {
		resp.Ref = defaultID
	}

	return resp, nil
}

func (lbf *llbBridgeForwarder) getImmutableRef(ctx context.Context, id string) (cache.ImmutableRef, error) {
	lbf.mu.Lock()
	ref, ok := lbf.refs[id]
	if !ok {
		if lbf.result != nil {
			if r, ok := lbf.result.FindRef(id); ok {
				ref = r
			}
		}
		if ref == nil {
			lbf.mu.Unlock()
			return nil, errors.Errorf("no such ref: %s, all %+v", id, maps.Keys(lbf.refs))
		}
	}
	lbf.mu.Unlock()
	if ref == nil {
		return nil, errors.Errorf("empty ref: %s", id)
	}

	r, err := ref.Result(ctx)
	if err != nil {
		return nil, lbf.wrapSolveError(err)
	}

	workerRef, ok := r.Sys().(*worker.WorkerRef)
	if !ok {
		return nil, errors.Errorf("invalid ref: %T", r.Sys())
	}

	return workerRef.ImmutableRef, nil
}

func (lbf *llbBridgeForwarder) ReadFile(ctx context.Context, req *pb.ReadFileRequest) (*pb.ReadFileResponse, error) {
	ctx = tracing.ContextWithSpanFromContext(ctx, lbf.callCtx)

	ref, err := lbf.getImmutableRef(ctx, req.Ref)
	if err != nil {
		return nil, err
	}

	newReq := cacheutil.ReadRequest{
		Filename: req.FilePath,
	}
	if r := req.Range; r != nil {
		newReq.Range = &cacheutil.FileRange{
			Offset: int(r.Offset),
			Length: int(r.Length),
		}
	}

	var m snapshot.Mountable
	if ref != nil {
		m, err = ref.Mount(ctx, true, session.NewGroup(lbf.sid))
		if err != nil {
			return nil, err
		}
	}

	dt, err := cacheutil.ReadFile(ctx, m, newReq)
	if err != nil {
		return nil, lbf.wrapSolveError(err)
	}

	return &pb.ReadFileResponse{Data: dt}, nil
}

func (lbf *llbBridgeForwarder) ReadDir(ctx context.Context, req *pb.ReadDirRequest) (*pb.ReadDirResponse, error) {
	ctx = tracing.ContextWithSpanFromContext(ctx, lbf.callCtx)

	ref, err := lbf.getImmutableRef(ctx, req.Ref)
	if err != nil {
		return nil, err
	}

	newReq := cacheutil.ReadDirRequest{
		Path:           req.DirPath,
		IncludePattern: req.IncludePattern,
	}
	var m snapshot.Mountable
	if ref != nil {
		m, err = ref.Mount(ctx, true, session.NewGroup(lbf.sid))
		if err != nil {
			return nil, err
		}
	}
	entries, err := cacheutil.ReadDir(ctx, m, newReq)
	if err != nil {
		return nil, lbf.wrapSolveError(err)
	}

	return &pb.ReadDirResponse{Entries: entries}, nil
}

func (lbf *llbBridgeForwarder) StatFile(ctx context.Context, req *pb.StatFileRequest) (*pb.StatFileResponse, error) {
	ctx = tracing.ContextWithSpanFromContext(ctx, lbf.callCtx)

	ref, err := lbf.getImmutableRef(ctx, req.Ref)
	if err != nil {
		return nil, err
	}
	var m snapshot.Mountable
	if ref != nil {
		m, err = ref.Mount(ctx, true, session.NewGroup(lbf.sid))
		if err != nil {
			return nil, err
		}
	}
	st, err := cacheutil.StatFile(ctx, m, req.Path)
	if err != nil {
		return nil, err
	}

	return &pb.StatFileResponse{Stat: st}, nil
}

func (lbf *llbBridgeForwarder) Evaluate(ctx context.Context, req *pb.EvaluateRequest) (*pb.EvaluateResponse, error) {
	ctx = tracing.ContextWithSpanFromContext(ctx, lbf.callCtx)

	_, err := lbf.getImmutableRef(ctx, req.Ref)
	if err != nil {
		return nil, err
	}
	return &pb.EvaluateResponse{}, nil
}

func (lbf *llbBridgeForwarder) Ping(context.Context, *pb.PingRequest) (*pb.PongResponse, error) {
	workers := lbf.workers.WorkerInfos()
	pbWorkers := make([]*apitypes.WorkerRecord, 0, len(workers))
	for _, w := range workers {
		pbWorkers = append(pbWorkers, &apitypes.WorkerRecord{
			ID:        w.ID,
			Labels:    w.Labels,
			Platforms: opspb.PlatformsFromSpec(w.Platforms),
		})
	}

	return &pb.PongResponse{
		FrontendAPICaps: pb.Caps.All(),
		Workers:         pbWorkers,
		LLBCaps:         opspb.Caps.All(),
	}, nil
}

func (lbf *llbBridgeForwarder) Return(ctx context.Context, in *pb.ReturnRequest) (*pb.ReturnResponse, error) {
	if in.Error != nil {
		return lbf.setResult(nil, grpcerrors.FromGRPC(status.ErrorProto(&spb.Status{
			Code:    in.Error.Code,
			Message: in.Error.Message,
			Details: in.Error.Details,
		})))
	}
	r := &frontend.Result{
		Metadata: in.Result.Metadata,
	}

	switch res := in.Result.Result.(type) {
	case *pb.Result_RefDeprecated:
		ref, err := lbf.cloneRef(res.RefDeprecated)
		if err != nil {
			return nil, err
		}
		r.SetRef(ref)
	case *pb.Result_RefsDeprecated:
		for k, id := range res.RefsDeprecated.Refs {
			ref, err := lbf.cloneRef(id)
			if err != nil {
				return nil, err
			}
			r.AddRef(k, ref)
		}
	case *pb.Result_Ref:
		ref, err := lbf.cloneRef(res.Ref.Id)
		if err != nil {
			return nil, err
		}
		r.SetRef(ref)
	case *pb.Result_Refs:
		for k, ref := range res.Refs.Refs {
			ref, err := lbf.cloneRef(ref.Id)
			if err != nil {
				return nil, err
			}
			r.AddRef(k, ref)
		}
	}

	if in.Result.Attestations != nil {
		for k, pbAtts := range in.Result.Attestations {
			for _, pbAtt := range pbAtts.Attestation {
				att, err := gwclient.AttestationFromPB[solver.ResultProxy](pbAtt)
				if err != nil {
					return nil, err
				}
				if pbAtt.Ref != nil {
					ref, err := lbf.cloneRef(pbAtt.Ref.Id)
					if err != nil {
						return nil, err
					}
					att.Ref = ref
				}
				r.AddAttestation(k, *att)
			}
		}
	}

	return lbf.setResult(r, nil)
}

func (lbf *llbBridgeForwarder) Inputs(ctx context.Context, in *pb.InputsRequest) (*pb.InputsResponse, error) {
	return &pb.InputsResponse{
		Definitions: lbf.inputs,
	}, nil
}

func (lbf *llbBridgeForwarder) NewContainer(ctx context.Context, in *pb.NewContainerRequest) (_ *pb.NewContainerResponse, err error) {
	bklog.G(ctx).Debugf("|<--- NewContainer %s", in.ContainerID)
	ctrReq := container.NewContainerRequest{
		ContainerID: in.ContainerID,
		NetMode:     in.Network,
		Hostname:    in.Hostname,
		Platform:    in.Platform,
		Constraints: in.Constraints,
	}

	for _, m := range in.Mounts {
		var workerRef *worker.WorkerRef
		if m.ResultID != "" {
			var ok bool
			workerRef, ok = lbf.workerRefByID[m.ResultID]
			if !ok {
				refProxy, err := lbf.convertRef(m.ResultID)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to find ref %s for %q mount", m.ResultID, m.Dest)
				}

				res, err := refProxy.Result(ctx)
				if err != nil {
					return nil, stack.Enable(err)
				}

				workerRef, ok = res.Sys().(*worker.WorkerRef)
				if !ok {
					return nil, errors.Errorf("invalid reference %T", res.Sys())
				}
			}
		}
		ctrReq.Mounts = append(ctrReq.Mounts, container.Mount{
			WorkerRef: workerRef,
			Mount: &opspb.Mount{
				Dest:      m.Dest,
				Selector:  m.Selector,
				Readonly:  m.Readonly,
				MountType: m.MountType,
				CacheOpt:  m.CacheOpt,
				SecretOpt: m.SecretOpt,
				SSHOpt:    m.SSHOpt,
			},
		})
	}

	// Not using `ctx` here because it will get cancelled as soon as NewContainer returns
	// and we want the context to live for the duration of the container.
	group := session.NewGroup(lbf.sid)

	cm, err := lbf.workers.DefaultCacheManager()
	if err != nil {
		return nil, stack.Enable(err)
	}

	ctrReq.ExtraHosts, err = container.ParseExtraHosts(in.ExtraHosts)
	if err != nil {
		return nil, stack.Enable(err)
	}

	ctr, err := container.NewContainer(context.Background(), cm, lbf.executor, lbf.sm, group, ctrReq)
	if err != nil {
		return nil, stack.Enable(err)
	}
	defer func() {
		if err != nil {
			ctr.Release(context.WithoutCancel(ctx)) // ensure release on error
		}
	}()

	lbf.ctrsMu.Lock()
	defer lbf.ctrsMu.Unlock()
	// ensure we are not clobbering a dup container id request
	if _, ok := lbf.ctrs[in.ContainerID]; ok {
		return nil, stack.Enable(status.Errorf(codes.AlreadyExists, "Container %s already exists", in.ContainerID))
	}
	lbf.ctrs[in.ContainerID] = ctr
	return &pb.NewContainerResponse{}, nil
}

func (lbf *llbBridgeForwarder) ReleaseContainer(ctx context.Context, in *pb.ReleaseContainerRequest) (*pb.ReleaseContainerResponse, error) {
	bklog.G(ctx).Debugf("|<--- ReleaseContainer %s", in.ContainerID)
	lbf.ctrsMu.Lock()
	ctr, ok := lbf.ctrs[in.ContainerID]
	delete(lbf.ctrs, in.ContainerID)
	lbf.ctrsMu.Unlock()
	if !ok {
		return nil, errors.Errorf("container details for %s not found", in.ContainerID)
	}
	err := ctr.Release(ctx)
	return &pb.ReleaseContainerResponse{}, stack.Enable(err)
}

func (lbf *llbBridgeForwarder) Warn(ctx context.Context, in *pb.WarnRequest) (*pb.WarnResponse, error) {
	// validate ranges are valid
	for _, r := range in.Ranges {
		if r == nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid source range")
		}
	}
	err := lbf.llbBridge.Warn(ctx, digest.Digest(in.Digest), string(in.Short), frontend.WarnOpts{
		Level:      int(in.Level),
		SourceInfo: in.Info,
		Range:      in.Ranges,
		Detail:     in.Detail,
		URL:        in.Url,
	})
	if err != nil {
		return nil, err
	}
	return &pb.WarnResponse{}, nil
}

type processIO struct {
	id       string
	mu       sync.Mutex
	resize   func(context.Context, gwclient.WinSize) error
	signal   func(context.Context, syscall.Signal) error
	done     chan struct{}
	doneOnce sync.Once
	// these track the process side of the io pipe for
	// read (fd=0) and write (fd=1, fd=2)
	processReaders map[uint32]io.ReadCloser
	processWriters map[uint32]io.WriteCloser
	// these track the server side of the io pipe, so
	// when we receive an EOF over grpc, we will close
	// this end
	serverWriters map[uint32]io.WriteCloser
	serverReaders map[uint32]io.ReadCloser
}

func newProcessIO(id string, openFds []uint32) *processIO {
	pio := &processIO{
		id:             id,
		processReaders: map[uint32]io.ReadCloser{},
		processWriters: map[uint32]io.WriteCloser{},
		serverReaders:  map[uint32]io.ReadCloser{},
		serverWriters:  map[uint32]io.WriteCloser{},
		done:           make(chan struct{}),
	}

	for _, fd := range openFds {
		// TODO do we know which way to pipe each fd?  For now assume fd0 is for
		// reading, and the rest are for writing
		r, w := io.Pipe()
		if fd == 0 {
			pio.processReaders[fd] = r
			pio.serverWriters[fd] = w
		} else {
			pio.processWriters[fd] = w
			pio.serverReaders[fd] = r
		}
	}

	return pio
}

func (pio *processIO) Close() (err error) {
	pio.mu.Lock()
	defer pio.mu.Unlock()
	for fd, r := range pio.processReaders {
		delete(pio.processReaders, fd)
		err1 := r.Close()
		if err1 != nil && err == nil {
			err = stack.Enable(err1)
		}
	}
	for fd, w := range pio.serverReaders {
		delete(pio.serverReaders, fd)
		err1 := w.Close()
		if err1 != nil && err == nil {
			err = stack.Enable(err1)
		}
	}
	pio.Done()
	return err
}

func (pio *processIO) Discard() {
	pio.mu.Lock()
	defer pio.mu.Unlock()
	pio.processReaders = nil
	pio.processWriters = nil
	pio.serverReaders = nil
	pio.serverWriters = nil
	pio.Done()
}

// hold lock before calling
func (pio *processIO) Done() {
	stillOpen := len(pio.processReaders) + len(pio.processWriters) + len(pio.serverReaders) + len(pio.serverWriters)
	if stillOpen == 0 {
		pio.doneOnce.Do(func() {
			close(pio.done)
		})
	}
}

func (pio *processIO) Write(f *pb.FdMessage) (err error) {
	pio.mu.Lock()
	writer := pio.serverWriters[f.Fd]
	pio.mu.Unlock()
	if writer == nil {
		return status.Errorf(codes.OutOfRange, "fd %d unavailable to write", f.Fd)
	}
	defer func() {
		if err != nil || f.EOF {
			writer.Close()
			pio.mu.Lock()
			defer pio.mu.Unlock()
			delete(pio.serverWriters, f.Fd)
			pio.Done()
		}
	}()
	if len(f.Data) > 0 {
		_, err = writer.Write(f.Data)
		return stack.Enable(err)
	}
	return nil
}

type outputWriter struct {
	stream    pb.LLBBridge_ExecProcessServer
	fd        uint32
	processID string
}

func (w *outputWriter) Write(msg []byte) (int, error) {
	bklog.G(w.stream.Context()).Debugf("|---> File Message %s, fd=%d, %d bytes", w.processID, w.fd, len(msg))
	err := w.stream.Send(&pb.ExecMessage{
		ProcessID: w.processID,
		Input: &pb.ExecMessage_File{
			File: &pb.FdMessage{
				Fd:   w.fd,
				Data: msg,
			},
		},
	})
	return len(msg), stack.Enable(err)
}

type execProcessServerThreadSafe struct {
	pb.LLBBridge_ExecProcessServer
	sendMu sync.Mutex
}

func (w *execProcessServerThreadSafe) Send(m *pb.ExecMessage) error {
	w.sendMu.Lock()
	defer w.sendMu.Unlock()
	return w.LLBBridge_ExecProcessServer.Send(m)
}

func (lbf *llbBridgeForwarder) ExecProcess(srv pb.LLBBridge_ExecProcessServer) error {
	srv = &execProcessServerThreadSafe{
		LLBBridge_ExecProcessServer: srv,
	}
	eg, ctx := errgroup.WithContext(srv.Context())

	msgs := make(chan *pb.ExecMessage)

	eg.Go(func() error {
		defer close(msgs)
		for {
			execMsg, err := srv.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return stack.Enable(err)
			}
			switch m := execMsg.GetInput().(type) {
			case *pb.ExecMessage_Init:
				bklog.G(ctx).Debugf("|<--- Init Message %s", execMsg.ProcessID)
			case *pb.ExecMessage_File:
				if m.File.EOF {
					bklog.G(ctx).Debugf("|<--- File Message %s, fd=%d, EOF", execMsg.ProcessID, m.File.Fd)
				} else {
					bklog.G(ctx).Debugf("|<--- File Message %s, fd=%d, %d bytes", execMsg.ProcessID, m.File.Fd, len(m.File.Data))
				}
			case *pb.ExecMessage_Resize:
				bklog.G(ctx).Debugf("|<--- Resize Message %s", execMsg.ProcessID)
			case *pb.ExecMessage_Signal:
				bklog.G(ctx).Debugf("|<--- Signal Message %s: %s", execMsg.ProcessID, m.Signal.Name)
			}
			select {
			case <-ctx.Done():
			case msgs <- execMsg:
			}
		}
	})

	eg.Go(func() error {
		pios := make(map[string]*processIO)
		// close any stray pios on exit to make sure
		// all the associated resources get cleaned up
		defer func() {
			for _, pio := range pios {
				pio.Close()
				pio.Discard()
			}
		}()

		for {
			var execMsg *pb.ExecMessage
			select {
			case <-ctx.Done():
				return nil
			case execMsg = <-msgs:
			}
			if execMsg == nil {
				return nil
			}

			pid := execMsg.ProcessID
			if pid == "" {
				return stack.Enable(status.Errorf(codes.InvalidArgument, "ProcessID required"))
			}

			pio, pioFound := pios[pid]

			if data := execMsg.GetFile(); data != nil {
				if !pioFound {
					return stack.Enable(status.Errorf(codes.NotFound, "IO for process %q not found", pid))
				}
				err := pio.Write(data)
				if err != nil {
					return stack.Enable(err)
				}
			} else if resize := execMsg.GetResize(); resize != nil {
				if !pioFound {
					return stack.Enable(status.Errorf(codes.NotFound, "IO for process %q not found", pid))
				}
				pio.resize(ctx, gwclient.WinSize{
					Cols: resize.Cols,
					Rows: resize.Rows,
				})
			} else if sig := execMsg.GetSignal(); sig != nil {
				if !pioFound {
					return stack.Enable(status.Errorf(codes.NotFound, "IO for process %q not found", pid))
				}
				syscallSignal, ok := signal.SignalMap[sig.Name]
				if !ok {
					return stack.Enable(status.Errorf(codes.InvalidArgument, "Unknown signal %s", sig.Name))
				}
				pio.signal(ctx, syscallSignal)
			} else if init := execMsg.GetInit(); init != nil {
				if pioFound {
					return stack.Enable(status.Errorf(codes.AlreadyExists, "Process %s already exists", pid))
				}
				id := init.ContainerID
				lbf.ctrsMu.Lock()
				ctr, ok := lbf.ctrs[id]
				lbf.ctrsMu.Unlock()
				if !ok {
					return stack.Enable(status.Errorf(codes.NotFound, "container %q previously released or not created", id))
				}

				initCtx, initCancel := context.WithCancelCause(context.Background())
				defer initCancel(errors.WithStack(context.Canceled))

				pio := newProcessIO(pid, init.Fds)
				pios[pid] = pio

				proc, err := ctr.Start(initCtx, gwclient.StartRequest{
					Args:                      init.Meta.Args,
					Env:                       init.Meta.Env,
					SecretEnv:                 init.Secretenv,
					User:                      init.Meta.User,
					Cwd:                       init.Meta.Cwd,
					Tty:                       init.Tty,
					Stdin:                     pio.processReaders[0],
					Stdout:                    pio.processWriters[1],
					Stderr:                    pio.processWriters[2],
					SecurityMode:              init.Security,
					RemoveMountStubsRecursive: init.Meta.RemoveMountStubsRecursive,
				})
				if err != nil {
					return stack.Enable(err)
				}
				pio.resize = proc.Resize
				pio.signal = proc.Signal

				eg.Go(func() error {
					<-pio.done
					bklog.G(ctx).Debugf("|---> Done Message %s", pid)
					err := srv.Send(&pb.ExecMessage{
						ProcessID: pid,
						Input: &pb.ExecMessage_Done{
							Done: &pb.DoneMessage{},
						},
					})
					return stack.Enable(err)
				})

				eg.Go(func() error {
					defer func() {
						pio.Close()
					}()
					err := proc.Wait()

					var statusCode uint32
					var exitError *pb.ExitError
					var statusError *spb.Status
					if err != nil {
						statusCode = pb.UnknownExitStatus
						st, _ := status.FromError(grpcerrors.ToGRPC(ctx, err))
						stp := st.Proto()
						statusError = &spb.Status{
							Code:    stp.Code,
							Message: stp.Message,
							Details: stp.Details,
						}
					}
					if errors.As(err, &exitError) {
						statusCode = exitError.ExitCode
					}
					bklog.G(ctx).Debugf("|---> Exit Message %s, code=%d, error=%s", pid, statusCode, err)
					sendErr := srv.Send(&pb.ExecMessage{
						ProcessID: pid,
						Input: &pb.ExecMessage_Exit{
							Exit: &pb.ExitMessage{
								Code:  statusCode,
								Error: statusError,
							},
						},
					})

					if sendErr != nil && err != nil {
						return errors.Wrap(sendErr, err.Error())
					} else if sendErr != nil {
						return stack.Enable(sendErr)
					}

					if err != nil && statusCode != 0 {
						// this was a container exit error which is "normal" so
						// don't return this error from the errgroup
						return nil
					}
					return stack.Enable(err)
				})

				bklog.G(ctx).Debugf("|---> Started Message %s", pid)
				err = srv.Send(&pb.ExecMessage{
					ProcessID: pid,
					Input: &pb.ExecMessage_Started{
						Started: &pb.StartedMessage{},
					},
				})
				if err != nil {
					return stack.Enable(err)
				}

				// start sending Fd output back to client, this is done after
				// StartedMessage so that Fd output will not potentially arrive
				// to the client before "Started" as the container starts up.
				for fd, file := range pio.serverReaders {
					fd, file := fd, file
					eg.Go(func() error {
						defer func() {
							file.Close()
							pio.mu.Lock()
							defer pio.mu.Unlock()
							w := pio.processWriters[fd]
							if w != nil {
								w.Close()
							}
							delete(pio.processWriters, fd)
							pio.Done()
						}()
						dest := &outputWriter{
							stream:    srv,
							fd:        fd,
							processID: pid,
						}
						_, err := io.Copy(dest, file)
						// ignore ErrClosedPipe, it is EOF for our usage.
						if err != nil && !errors.Is(err, io.ErrClosedPipe) {
							return stack.Enable(err)
						}
						// no error so must be EOF
						bklog.G(ctx).Debugf("|---> File Message %s, fd=%d, EOF", pid, fd)
						err = srv.Send(&pb.ExecMessage{
							ProcessID: pid,
							Input: &pb.ExecMessage_File{
								File: &pb.FdMessage{
									Fd:  fd,
									EOF: true,
								},
							},
						})
						return stack.Enable(err)
					})
				}
			}
		}
	})

	err := eg.Wait()
	return stack.Enable(err)
}

func (lbf *llbBridgeForwarder) convertRef(id string) (solver.ResultProxy, error) {
	if id == "" {
		return nil, nil
	}

	lbf.mu.Lock()
	defer lbf.mu.Unlock()

	r, ok := lbf.refs[id]
	if !ok {
		return nil, errors.Errorf("return reference %s not found", id)
	}
	return r, nil
}

func (lbf *llbBridgeForwarder) cloneRef(id string) (solver.ResultProxy, error) {
	if id == "" {
		return nil, nil
	}

	lbf.mu.Lock()
	defer lbf.mu.Unlock()

	r, ok := lbf.refs[id]
	if !ok {
		return nil, errors.Errorf("return reference %s not found", id)
	}

	s1, s2 := solver.SplitResultProxy(r)
	lbf.refs[id] = s1
	return s2, nil
}

func serve(ctx context.Context, grpcServer *grpc.Server, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	bklog.G(ctx).Debugf("serving grpc connection")
	(&http2.Server{}).ServeConn(conn, &http2.ServeConnOpts{Handler: grpcServer})
}

type markTypeFrontend struct{}

func (*markTypeFrontend) SetImageOption(ii *llb.ImageInfo) {
	ii.RecordType = string(client.UsageRecordTypeFrontend)
}

func getCaps(label string) map[string]struct{} {
	if label == "" {
		return make(map[string]struct{})
	}
	caps := strings.Split(label, ",")
	out := make(map[string]struct{}, len(caps))
	for _, c := range caps {
		name := strings.SplitN(c, "+", 2)
		if name[0] != "" {
			out[name[0]] = struct{}{}
		}
	}
	return out
}

func addCapsForKnownFrontends(caps map[string]struct{}, dgst digest.Digest) {
	// these frontends were built without caps detection but do support inputs
	defaults := map[digest.Digest]struct{}{
		"sha256:9ac1c43a60e31dca741a6fe8314130a9cd4c4db0311fbbc636ff992ef60ae76d": {}, // docker/dockerfile:1.1.6
		"sha256:080bd74d8778f83e7b670de193362d8c593c8b14f5c8fb919d28ee8feda0d069": {}, // docker/dockerfile:1.1.7
		"sha256:60543a9d92b92af5088fb2938fb09b2072684af8384399e153e137fe081f8ab4": {}, // docker/dockerfile:1.1.6-experimental
		"sha256:de85b2f3a3e8a2f7fe48e8e84a65f6fdd5cd5183afa6412fff9caa6871649c44": {}, // docker/dockerfile:1.1.7-experimental
	}
	if _, ok := defaults[dgst]; ok {
		caps["moby.buildkit.frontend.inputs"] = struct{}{}
	}
}
