package file

import (
	"path/filepath"
	"reflect"
	"testing"

	"go.lsp.dev/protocol"
)

func TestURI(t *testing.T) {
	typ := reflect.TypeOf(RelativeURI(""))
	t.Log(typ.Name(), typ.PkgPath())

	urityp := reflect.TypeOf(protocol.URI(""))
	t.Log(urityp.Name(), urityp.PkgPath())

	t.Log(filepath.Rel("/a/b/c", "file:///a/b/c/d.txt"))
	t.Log(filepath.Rel("/a/b/c", "/a/b/c/d.txt"))

	t.Log(filepath.IsAbs("file:///a/b/c.txt"))
	t.Log(filepath.IsAbs(TrimScheme("file:///a/b/c.txt")))
}
