package file

import (
	"context"

	"github.com/bufbuild/buf/private/pkg/storage"
	"go.lsp.dev/protocol"
)

type Handle interface {
	URI() protocol.DocumentURI

	Version() int32

	Content() ([]byte, error)
}

type HasPath interface {
	Path() RelativeURI
}

type HandleWithPath interface {
	Handle
	HasPath
}

// A Source maps URIs to Handles.
type Source interface {

	// ReadFile will get a file handle with `DocumentURI` or `RelativeURI`
	ReadFile(ctx context.Context, uri protocol.URI) (Handle, error)

	// Close will remove a file handle cache with `DocumentURI` or `RelativeURI`
	Close(uri protocol.URI)

	// Stat will get a file ObjectInfo with `DocumentURI` or `RelativeURI`
	Stat(ctx context.Context, uri protocol.URI) (storage.ObjectInfo, error)
}
