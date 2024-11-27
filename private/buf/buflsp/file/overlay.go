package file

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/bufbuild/buf/private/pkg/storage"
	"go.lsp.dev/protocol"
)

type Overlay struct {
	root    string
	uri     protocol.DocumentURI
	content []byte
	version int32
}

func NewOverlay(uri protocol.DocumentURI, content []byte, version int32) *Overlay {
	return &Overlay{
		uri:     uri,
		content: content,
		version: version,
	}
}

func (o *Overlay) Change(content []byte, version int32) {
	o.content = content
	o.version = version
}

func (o *Overlay) URI() protocol.DocumentURI { return o.uri }

func (o *Overlay) Content() ([]byte, error) { return o.content, nil }
func (o *Overlay) Version() int32           { return o.version }

func (o *Overlay) Path() string {
	// open has checked this err
	rel, _ := filepath.Rel(o.root, o.uri.Filename())
	return rel
}

func (o *Overlay) ExternalPath() string {
	return ""
}

func (o *Overlay) LocalPath() string {
	return o.URI().Filename()
}

func (o *Overlay) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("URI", string(o.URI())),
		slog.String("Path", o.Path()),
		slog.Any("Version", o.Version()),
	)
}

// An OverlayFS is a file.Source that keeps track of overlays on top of a
// delegate FileSource.
type OverlayFS struct {
	root     string
	delegate Source

	mu       sync.Mutex
	overlays map[protocol.DocumentURI]*Overlay
}

func NewOverlayFS(root string, delegate Source) *OverlayFS {
	return &OverlayFS{
		root:     root,
		delegate: delegate,
		overlays: make(map[protocol.DocumentURI]*Overlay),
	}
}

func (fs *OverlayFS) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Type", "OverlayFS"),
		slog.String("Root", fs.root),
		slog.Int("Count", len(fs.overlays)),
	)
}

// Overlays returns a new unordered array of overlays.
func (fs *OverlayFS) Overlays() []*Overlay {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	overlays := make([]*Overlay, 0, len(fs.overlays))
	for _, overlay := range fs.overlays {
		overlays = append(overlays, overlay)
	}
	return overlays
}

func (fs *OverlayFS) ReadFile(ctx context.Context, uri protocol.URI) (Handle, error) {
	uri = fs.checkURI(uri)

	fs.mu.Lock()
	overlay, ok := fs.overlays[uri]
	fs.mu.Unlock()
	if ok {
		return overlay, nil
	}
	return fs.delegate.ReadFile(ctx, uri)
}

func (fs *OverlayFS) Get(uri protocol.DocumentURI) (*Overlay, error) {
	if uri == "" {
		return nil, errors.New("empty DocumentURI")
	}

	uri = fs.checkURI(uri)

	fs.mu.Lock()
	defer fs.mu.Unlock()
	if oo, ok := fs.overlays[uri]; ok {
		return oo, nil
	} else {
		return nil, os.ErrNotExist
	}
}

func (fs *OverlayFS) Open(o *Overlay) (*Overlay, error) {
	if o.uri == "" {
		return nil, errors.New("empty DocumentURI")
	}

	if IsDocumentURI(o.uri) {
		if _, err := filepath.Rel(fs.root, o.uri.Filename()); err != nil {
			return nil, errors.New("overlay fs not include this file")
		}
	}

	o.uri = fs.checkURI(o.uri)
	// fill root
	o.root = fs.root

	fs.mu.Lock()
	defer fs.mu.Unlock()
	if oo, ok := fs.overlays[o.uri]; ok {
		oo.Change(o.content, o.version)
		return oo, nil
	}

	fs.overlays[o.uri] = o
	return o, nil
}

func (fs *OverlayFS) Close(uri protocol.URI) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	uri = fs.checkURI(uri)

	delete(fs.overlays, uri)
}

func (fs *OverlayFS) Stat(ctx context.Context, uri protocol.DocumentURI) (storage.ObjectInfo, error) {
	uri = fs.checkURI(uri)

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if ol, ok := fs.overlays[uri]; ok {
		return ol, nil
	} else {
		return nil, os.ErrNotExist
	}
}

func (fs *OverlayFS) checkURI(u protocol.URI) protocol.DocumentURI {
	if !IsDocumentURI(u) {
		// as RelativeURI
		u = NormalURIStr(filepath.Join(fs.root, string(u)))
	}

	return u
}
