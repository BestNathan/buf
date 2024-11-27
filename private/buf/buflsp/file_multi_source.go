package buflsp

import (
	"context"
	"log/slog"
	"os"

	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/buf/private/pkg/storage"
	"go.lsp.dev/protocol"
)

type MultiSource struct {
	logger  *slog.Logger
	sources []file.Source
}

func NewMultiSource(logger *slog.Logger, sources ...file.Source) *MultiSource {

	ss := []file.Source{}

	for _, src := range sources {
		if ms, ok := src.(*MultiSource); ok {
			ss = append(ss, ms.Sources()...)
		} else {
			ss = append(ss, src)
		}
	}

	return &MultiSource{
		logger:  logger,
		sources: ss,
	}
}

func (ms *MultiSource) Sources() []file.Source {
	return ms.sources
}

func (m *MultiSource) ReadFile(ctx context.Context, uri protocol.DocumentURI) (file.Handle, error) {
	for _, s := range m.sources {
		if h, err := s.ReadFile(ctx, uri); err != nil {
			// m.logger.Debug(
			// 	"multi source read file fail",
			// 	"URI", uri,
			// 	"Source", s,
			// 	"Error", err,
			// )
		} else {
			// m.logger.Debug(
			// 	"multi source read file success",
			// 	"URI", h.URI(),
			// 	"Version", h.Version(),
			// 	"Source", s,
			// )

			return h, nil
		}
	}

	return nil, os.ErrNotExist
}

func (m *MultiSource) Close(uri protocol.DocumentURI) {
	for _, s := range m.sources {
		s.Close(uri)
	}
}

func (m *MultiSource) Stat(ctx context.Context, uri protocol.DocumentURI) (storage.ObjectInfo, error) {
	for _, s := range m.sources {
		if info, err := s.Stat(ctx, uri); err != nil {
			// m.logger.Debug(
			// 	"multi source stat file fail",
			// 	"URI", uri,
			// 	"Source", s,
			// 	"Error", err,
			// )
		} else {
			return info, nil
		}
	}

	return nil, os.ErrNotExist
}
