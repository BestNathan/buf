package parser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/bufbuild/buf/private/buf/bufformat"
	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/linker"
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/bufbuild/protocompile/walk"
	"go.lsp.dev/protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrUnParseFile = errors.New("file not parsed")
)

const (
	FileInitVersion int32 = -2
	descriptorPath        = "google/protobuf/descriptor.proto"
)

type File struct {
	handle   file.Handle
	reporter reporter.Reporter

	// compile
	parsemu    sync.RWMutex
	resolver   protocompile.Resolver
	symbols    *linker.Symbols
	linkerfile linker.File

	// different from file.Handle.URI
	// this uri for found imports
	// like: xxx/xxx/xxx.proto
	uri     file.RelativeURI
	content []byte
	version int32

	fn      *ast.FileNode
	res     parser.Result
	pkg     *ast.PackageNode
	imports []*ast.ImportNode
}

func newFile(
	uri file.RelativeURI,
	handle file.Handle,
	rpt reporter.Reporter,
	resolver protocompile.Resolver,
) *File {
	if strings.Contains(file.URI2Filename(uri), descriptorPath) {
		rpt = noopReporter
	}

	return &File{
		uri:      uri,
		handle:   handle,
		reporter: rpt,
		resolver: resolver,
		version:  FileInitVersion,
	}
}

func (f *File) reset(h file.Handle) {
	f.parsemu.Lock()
	defer f.parsemu.Unlock()

	if h != nil {
		f.handle = h
	}

	// reset handle info
	f.content = nil
	f.version = FileInitVersion

	// reset parse result
	f.fn = nil
	f.res = nil
	f.pkg = nil
	f.imports = nil
}

func (f *File) LookUpSymbol(name protoreflect.FullName) ast.SourceSpan {
	if f.symbols == nil {
		return nil
	}
	return f.symbols.Lookup(name)
}

func (f *File) Range() protocol.Range {
	if f.fn == nil {
		return protocol.Range{}
	} else {
		ni := f.NodeInfo(f.fn)

		return protocol.Range{
			Start: NewZeroBaseSourcePos(ni.Start()).ToPosition(),
			End:   NewZeroBaseSourcePos(ni.End()).ToPosition(),
		}
	}
}

func (f *File) Reset() {
	f.reset(nil)
}

func (f *File) ResetWithHandle(h file.Handle) {
	f.reset(h)
}

func (f *File) Format(out io.Writer) error {
	f.parsemu.RLock()
	defer f.parsemu.RUnlock()

	if f.fn == nil {
		return ErrUnParseFile
	}

	return bufformat.FormatFileNode(out, f.fn)
}

func (f *File) Walk(visitor ast.Visitor) error {
	if f.fn == nil {
		return ErrUnParseFile
	}
	f.parsemu.RLock()
	defer f.parsemu.RUnlock()

	return ast.Walk(f.fn, visitor)
}

func (f *File) WalkDescriptorProto(fn func(protoreflect.FullName, proto.Message) error) error {
	if f.res == nil {
		return ErrUnParseFile
	}
	f.parsemu.RLock()
	defer f.parsemu.RUnlock()

	return walk.DescriptorProtos(f.res.FileDescriptorProto(), fn)
}

func (f *File) NodeInfo(n ast.Node) ast.NodeInfo {
	return f.fn.NodeInfo(n)
}

func (f *File) needReParse() bool {
	if f.fn == nil {
		return true
	}

	hv := f.handle.Version()

	if f.version < hv {
		return true
	}

	return false
}

func (f *File) Parse() error {
	if f.handle == nil {
		return errors.New("nil file handle")
	}

	f.parsemu.Lock()
	defer f.parsemu.Unlock()

	if !f.needReParse() {
		return nil
	}

	content, err := f.handle.Content()
	if err != nil {
		return fmt.Errorf("file content: %w", err)
	}

	// parse
	fn, err := parser.Parse(
		file.URI2Filename(f.uri),
		bytes.NewReader(content),
		reporter.NewHandler(f.reporter),
	)
	if err != nil && !errors.Is(err, reporter.ErrInvalidSource) {
		return fmt.Errorf("proto compile parse: %w", err)
	}

	// validate simple
	res, err := parser.ResultFromAST(
		fn,
		true,
		reporter.NewHandler(f.reporter),
	)
	if err != nil && !errors.Is(err, reporter.ErrInvalidSource) {
		return fmt.Errorf("proto compile result from ast: %w", err)
	}

	// handle info
	f.content = content
	f.version = f.handle.Version()

	// parse result
	f.fn = fn
	f.res = res

	// future info storage may be included here
	_ = ast.Walk(fn, &ast.SimpleVisitor{
		DoVisitPackageNode: func(pn *ast.PackageNode) error {
			// pkg
			f.pkg = pn
			return nil
		},
		DoVisitImportNode: func(in *ast.ImportNode) error {
			// imports
			f.imports = append(f.imports, in)
			return nil
		},
	})

	return nil
}

func (f *File) Compile() error {
	f.parsemu.Lock()
	defer f.parsemu.Unlock()

	f.symbols = &linker.Symbols{}

	compiler := &protocompile.Compiler{
		Resolver:       f.resolver,
		Reporter:       f.reporter,
		SourceInfoMode: protocompile.SourceInfoExtraOptionLocations,
		Symbols:        f.symbols,
	}

	compiled, err := compiler.Compile(context.Background(), file.URI2Filename(f.uri))
	if err != nil && !errors.Is(err, reporter.ErrInvalidSource) {
		return fmt.Errorf("compile: %w", err)
	}

	if compiled[0] == nil {
		return errors.New("no compiled files")
	}

	f.linkerfile = compiled[0]
	return nil
}

func (f *File) PackageName() PackageName {
	f.parsemu.RLock()
	defer f.parsemu.RUnlock()

	if f.pkg == nil {
		return ""
	}

	return PackageName(f.pkg.Name.AsIdentifier())
}

func (f *File) Imports() []file.RelativeURI {
	f.parsemu.RLock()
	defer f.parsemu.RUnlock()

	is := []protocol.URI{}
	for _, im := range f.imports {
		is = append(is, file.NormalURIStr(im.Name.AsString()))
	}

	return is
}

func (f *File) URI() protocol.DocumentURI {
	return f.uri
}

func (f *File) Version() int32 {
	return f.version
}

func (f *File) Content() ([]byte, error) {
	return f.content, nil
}

// NodesAt find all nodes that will contain `pos.Line`
// and for node that located at one line, will test for `pos.Character`
func (f *File) NodesAt(pos protocol.Position) Nodes {
	ns := Nodes{}

	_ = f.Walk(&ast.SimpleVisitor{
		DoVisitNode: func(n ast.Node) error {
			ni := f.NodeInfo(n)
			if PositionWithinNode(ni, pos) {
				ns = append(ns, n)
			}

			return nil
		},
	})

	return ns
}

func (f *File) Nearest(pos protocol.Position) ast.Node {
	fn := f.fn
	if fn == nil {
		return nil
	}

	ns := f.NodesAt(pos)
	if len(ns) == 0 {
		return nil
	}

	var nearest ast.Node
	nearestline, nearestcol := -1, -1
	for _, n := range ns {
		func() {
			var assign bool
			ni := fn.NodeInfo(n)

			if !PositionWithinNode(ni, pos) {
				return
			}

			linelen, collen := PositionNodeLen(ni, pos)
			// line first then col
			// line same then col

			defer func() {
				if assign {
					nearest = n
					nearestline = linelen
					nearestcol = collen
				}
			}()

			// first node, just assign
			if nearest == nil {
				assign = true
				return
			}

			// cmp line
			if linelen > nearestline { // line gt, drop
				return
			} else if linelen < nearestline { // line lt, assign
				assign = true
			} else { // line eq
				if collen < nearestcol {
					assign = true
				}
			}

		}()
	}

	return nearest
}
