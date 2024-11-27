package parser

import (
	"strings"
	"sync"

	"github.com/bufbuild/protocompile/ast"
	"go.lsp.dev/protocol"
)

type Package struct {
	name   PackageName
	filemu sync.RWMutex
	files  map[protocol.DocumentURI]*File
}

func NewPackage(n PackageName) *Package {
	return &Package{
		name:  n,
		files: map[protocol.URI]*File{},
	}
}

func (p *Package) foreach(fn func(*File) (ifbreak bool)) {
	for _, f := range p.files {
		if fn(f) {
			break
		}
	}
}

func (p *Package) Walk(visitor ast.Visitor) (err error) {
	p.filemu.RLock()
	defer p.filemu.RUnlock()

	p.foreach(func(f *File) (ifbreak bool) {
		err = f.Walk(visitor)
		if err != nil {
			return true
		} else {
			return false
		}
	})

	return
}

func (p *Package) AddFile(f *File) {
	if f == nil {
		return
	}

	p.filemu.Lock()
	defer p.filemu.Unlock()
	if _, ok := p.files[f.URI()]; !ok {
		p.files[f.URI()] = f
	}
}

func (p *Package) RemoveFile(f *File) {
	if f == nil {
		return
	}

	p.filemu.Lock()
	defer p.filemu.Unlock()
	delete(p.files, f.URI())
}

type PackageName string

func (p PackageName) Parts() []string {
	return strings.Split(string(p), ".")
}
