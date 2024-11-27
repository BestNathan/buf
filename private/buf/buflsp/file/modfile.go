package file

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"context"

	"log/slog"

	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/pkg/storage"
	"go.lsp.dev/protocol"
)

type modfile struct {
	fileinfo bufmodule.FileInfo

	simplefile
}

func (f *modfile) Path() RelativeURI {
	if f.fileinfo == nil {
		return ""
	}
	return RelativeURI(f.fileinfo.Path())
}

func (m *modfile) LogValue() slog.Value {
	if m.err != nil {
		return slog.GroupValue(
			slog.Any("URI", m.uri),
			slog.Any("Error", m.err),
		)
	} else {
		return slog.GroupValue(
			slog.Any("URI", m.uri),
			slog.Any("Path", m.Path()),
			slog.Any("Length", len(m.content)),
		)
	}
}

type modSource struct {
	module bufmodule.Module

	mu       sync.RWMutex
	urifiles map[RelativeURI]*modfile
	relfiles map[RelativeURI]*modfile
}

func NewModSource(mod bufmodule.Module) Source {
	return &modSource{
		module:   mod,
		urifiles: make(map[protocol.DocumentURI]*modfile),
		relfiles: make(map[RelativeURI]*modfile),
	}
}

func (s *modSource) trylocal(uri protocol.URI) *modfile {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if f, ok := s.urifiles[uri]; ok {
		return f
	} else if f, ok := s.relfiles[uri]; ok {
		return f
	} else {
		return nil
	}
}

func (s *modSource) ReadFile(ctx context.Context, uri protocol.URI) (Handle, error) {
	if f := s.trylocal(uri); f != nil {
		return f, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	stat, err := s.internalstat(ctx, uri)
	if err != nil {
		return nil, err
	}

	mf := &modfile{
		fileinfo:   stat,
		simplefile: simplefile{uri: uri},
	}

	f, err := s.module.GetFile(ctx, stat.Path())
	if err != nil {
		mf.err = err
	} else {
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			mf.err = err
		} else {
			mf.content = b
			mf.fileinfo = f
		}
	}

	s.urifiles[mf.URI()] = mf
	s.relfiles[mf.Path()] = mf
	return mf, nil
}

func (s *modSource) closefile(f *modfile) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if f != nil {
		delete(s.urifiles, f.uri)
		delete(s.relfiles, f.Path())
	}
}

func (s *modSource) Close(uri protocol.URI) {
	if f := s.trylocal(uri); f != nil {
		s.closefile(f)
	}
}

func (s *modSource) Stat(ctx context.Context, uri protocol.URI) (storage.ObjectInfo, error) {
	if f := s.trylocal(uri); f != nil {
		return f.fileinfo, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.internalstat(ctx, uri)
}

func (s *modSource) internalstat(ctx context.Context, uri protocol.URI) (bufmodule.FileInfo, error) {
	uri = NormalURI(uri)
	isdoc := IsDocumentURI(uri)
	filename := URI2Filename(uri)

	var fi bufmodule.FileInfo
	err := s.module.WalkFileInfos(ctx, func(fileinfo bufmodule.FileInfo) error {
		// not proto
		if fileinfo.FileType() != bufmodule.FileTypeProto {
			return nil
		}

		if isdoc {
			if strings.EqualFold(filename, fileinfo.LocalPath()) {
				fi = fileinfo
			}
		} else {
			if strings.EqualFold(filename, fileinfo.Path()) {
				fi = fileinfo
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("module walk: %w", err)
	}

	if fi == nil {
		return nil, fmt.Errorf("stat file uri: `%s`: %w", uri, os.ErrNotExist)
	} else {
		return fi, nil
	}
}

func (m *modSource) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Type", "Module"),
		slog.Any("IsLocal", m.module.IsLocal()),
		slog.String("BucketID", m.module.BucketID()),
		slog.String("OpaqueID", m.module.OpaqueID()),
		slog.Any("FullName", m.module.ModuleFullName()),
		slog.String("Description", m.module.Description()),
	)
}
