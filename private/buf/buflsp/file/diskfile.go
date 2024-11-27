package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"log/slog"

	"github.com/bufbuild/buf/private/buf/buflsp/bucket"
	"github.com/bufbuild/buf/private/pkg/storage"
	"go.lsp.dev/protocol"
)

type diskFile struct {
	objectinfo storage.ObjectInfo
	simplefile
}

func (f *diskFile) Path() RelativeURI {
	if f.objectinfo == nil {
		return ""
	}
	return RelativeURI(f.objectinfo.Path())
}

func (f *diskFile) LogValue() slog.Value {
	if f.err != nil {
		return slog.GroupValue(
			slog.Any("URI", f.uri),
			slog.Any("Error", f.err),
		)
	} else {
		return slog.GroupValue(
			slog.Any("URI", f.uri),
			slog.Any("Path", f.Path()),
			slog.Any("Length", len(f.content)),
		)
	}
}

type diskSource struct {
	bucket storage.ReadWriteBucket

	mu    sync.Mutex
	files map[protocol.DocumentURI]*diskFile
}

func NewDiskSource() Source {
	return &diskSource{
		bucket: bucket.RootBucket(),
		files:  map[protocol.DocumentURI]*diskFile{},
	}
}

func (d *diskSource) Close(uri protocol.DocumentURI) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.files, uri)
}

func (d *diskSource) Stat(ctx context.Context, uri protocol.DocumentURI) (storage.ObjectInfo, error) {
	if !IsDocumentURI(uri) {
		return nil, ErrorInvalidURI(uri, ErrInvalidURIReasonOnlyAcceptDocumentURI)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if df, ok := d.files[uri]; ok {
		return df.objectinfo, nil
	}

	return d.bucket.Stat(ctx, uri.Filename())
}

func (d *diskSource) ReadFile(ctx context.Context, uri protocol.DocumentURI) (Handle, error) {
	if !IsDocumentURI(uri) {
		return nil, ErrorInvalidURI(uri, ErrInvalidURIReasonOnlyAcceptDocumentURI)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if df, ok := d.files[uri]; ok {
		return df, nil
	}

	// for root filename is path
	p := uri.Filename()

	stat, err := d.bucket.Stat(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("stat file uri: `%s`: %w", uri, os.ErrNotExist)
	}

	df := &diskFile{
		objectinfo: stat,
		simplefile: simplefile{uri: uri},
	}

	f, err := d.bucket.Get(ctx, p)
	if err != nil {
		df.err = err
	} else {
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			df.err = err
		} else {
			df.content = b
			df.objectinfo = f
		}
	}

	d.files[uri] = df
	return df, nil
}

func (d *diskSource) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Type", "Disk"),
		slog.Int("Count", len(d.files)),
	)
}
