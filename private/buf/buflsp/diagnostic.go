package buflsp

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"go.lsp.dev/protocol"
)

// - When a file changes it is the server's responsibility to re-compute diagnostics and push them to the client.
// - If the computed set is empty it has to push the empty array to clear former diagnostics.
// Newly pushed diagnostics always replace previously pushed diagnostics. There is no merging that happens on the client side.
type DiagnosticClient struct {
	logger *slog.Logger
	client protocol.Client

	mu    sync.Mutex
	cache map[protocol.DocumentURI]map[protocol.Range]protocol.Diagnostic
}

func NewDiagnosticClient(logger *slog.Logger, client protocol.Client) *DiagnosticClient {
	return &DiagnosticClient{
		logger: logger,
		client: client,
		cache:  map[protocol.URI]map[protocol.Range]protocol.Diagnostic{},
	}
}

func (d *DiagnosticClient) AddDiagnostics(f file.Handle, ds []protocol.Diagnostic) error {
	d.logger.Debug(
		"add diagnostics",
		"URI", f.URI(),
		"Diagnostics", ds,
	)

	uri := f.URI()

	d.mu.Lock()
	defer d.mu.Unlock()

	dmap, ok := d.cache[uri]
	if !ok {
		dmap = map[protocol.Range]protocol.Diagnostic{}
	}

	for _, d := range ds {
		dmap[d.Range] = d
	}

	d.cache[uri] = dmap

	return nil
}

func (d *DiagnosticClient) Notify(f file.Handle) error {
	uri := f.URI()

	d.mu.Lock()
	dmap := d.cache[uri]

	d.logger.Debug(
		"notify diagnostics",
		"URI", f.URI(),
		"Count", len(dmap),
	)

	ds := []protocol.Diagnostic{}
	for _, v := range dmap {
		ds = append(ds, v)
	}

	d.mu.Unlock()

	if err := d.client.PublishDiagnostics(context.Background(), &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: ds,
	}); err != nil {
		return fmt.Errorf("client publish diagnostics: %w", err)
	}

	return nil
}

func (d *DiagnosticClient) Reset(f file.Handle) {
	uri := f.URI()

	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.cache, uri)
}
