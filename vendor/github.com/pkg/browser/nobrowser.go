// Package nobrowser stubs out github.com/pkg/browser
//
// See https://github.com/moby/buildkit/pull/3010#issuecomment-1226847222
package nobrowser

import (
	"fmt"
	"io"
)

// Stdout is the io.Writer to which executed commands write standard output.
var Stdout io.Writer

// Stderr is the io.Writer to which executed commands write standard error.
var Stderr io.Writer

// OpenFile opens new browser window for the file path.
func OpenFile(path string) error {
	return fmt.Errorf("not implemented")
}

// OpenReader consumes the contents of r and presents the
// results in a new browser window.
func OpenReader(r io.Reader) error {
	return fmt.Errorf("not implemented")
}

// OpenURL opens a new browser window pointing to url.
func OpenURL(url string) error {
	return fmt.Errorf("not implemented")
}
