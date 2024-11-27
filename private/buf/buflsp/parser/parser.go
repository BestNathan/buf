package parser

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/protocompile"
)

type Parser struct {
	logger              *slog.Logger
	diagnosticCollector DiagnosticCollector
	resolver            protocompile.Resolver

	mu    sync.Mutex
	files map[file.RelativeURI]*File
	pkgs  map[PackageName]*Package
}

func NewParser(
	logger *slog.Logger,
	diagnosticCollector DiagnosticCollector,
	resolver protocompile.Resolver,
) *Parser {
	if diagnosticCollector == nil {
		diagnosticCollector = &noopDiagnosticCollector{}
	}

	return &Parser{
		logger:              logger,
		resolver:            resolver,
		diagnosticCollector: diagnosticCollector,
		files:               map[file.RelativeURI]*File{},
		pkgs:                map[PackageName]*Package{},
	}
}

func (p *Parser) Parse(uri file.RelativeURI, h file.Handle) (*File, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	f, ok := p.files[uri]
	if ok {
		f.ResetWithHandle(h)
		p.pkgremove(f)
	} else {
		f = newFile(
			uri,
			h,
			&fileReporter{
				logger:              p.logger,
				handle:              h,
				diagnosticCollector: p.diagnosticCollector,
			},
			p.resolver,
		)

		p.files[uri] = f
	}

	if err := f.Parse(); err != nil {
		return nil, fmt.Errorf("file parse: %w", err)
	}

	if err := f.Compile(); err != nil {
		p.logger.Debug(
			"compile fail",
			"URI", f.uri,
			"Error", err,
		)
	} else {
		p.logger.Debug(
			"compile success",
			"URI", f.uri,
		)
	}

	p.pkgadd(f)

	return f, nil
}

func (p *Parser) FindFilePkg(f *File) (*Package, bool) {
	if f == nil {
		return nil, false
	}
	pkgn := f.PackageName()

	if pkg, ok := p.pkgs[pkgn]; ok {
		return pkg, true
	} else {
		return nil, false
	}
}

func (p *Parser) pkgremove(f *File) {
	if f == nil {
		return
	}

	pkgn := f.PackageName()
	pkg, ok := p.pkgs[pkgn]
	if ok {
		pkg.RemoveFile(f)
	}
}

func (p *Parser) pkgadd(f *File) {
	if f == nil {
		return
	}

	pkgn := f.PackageName()
	pkg, ok := p.pkgs[pkgn]
	if !ok {
		pkg = NewPackage(pkgn)
		p.pkgs[pkgn] = pkg
	}

	pkg.AddFile(f)
}
