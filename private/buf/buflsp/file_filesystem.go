package buflsp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/bufbuild/buf/private/buf/bufctl"
	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/protocompile"
	"go.lsp.dev/protocol"
)

type filesystem struct {
	logger     *slog.Logger
	controller bufctl.Controller
	folders    []*folder
}

func NewFileSystem(logger *slog.Logger, ctlr bufctl.Controller) *filesystem {
	return &filesystem{logger: logger, controller: ctlr}
}

func (fs *filesystem) Source() file.Source {
	ss := []file.Source{}
	for _, l := range fs.folders {
		ss = append(ss, l.Souce())
	}
	return NewMultiSource(fs.logger, ss...)
}

func (fs *filesystem) Resolver() protocompile.Resolver {
	return &protocompile.SourceResolver{
		Accessor: func(path string) (io.ReadCloser, error) {
			handle, err := fs.Source().ReadFile(context.Background(), file.NormalURIStr(path))
			if err != nil {
				return nil, fmt.Errorf("fs source read file: %w", err)
			}

			content, err := handle.Content()
			if err != nil {
				return nil, fmt.Errorf("handle content: %w", err)
			}

			return io.NopCloser(bytes.NewReader(content)), nil
		},
	}
}

func (fs *filesystem) Open(overlay *file.Overlay) (file.Handle, error) {
	for _, f := range fs.folders {
		if h, err := f.overlayfs.Open(overlay); err != nil {
			fs.logger.Debug(
				"open overlay fail",
				"OverlayFS", f.overlayfs,
				"Overlay", overlay,
				"Error", err,
			)
		} else {
			fs.logger.Debug(
				"open overlay success",
				"OverlayFS", f.overlayfs,
				"Overlay", overlay,
				"Handle", h,
			)

			return h, nil
		}
	}

	return nil, os.ErrNotExist
}

func (fs *filesystem) Location(ctx context.Context, uri file.RelativeURI) []protocol.Location {
	locs := []protocol.Location{}
	for _, f := range fs.folders {
		locs = append(locs, f.Location(ctx, uri)...)
	}
	return locs
}

func (fs *filesystem) init(ctx context.Context, folders protocol.WorkspaceFolders) error {
	if fs.controller == nil {
		return errors.New("nil controller")
	}

	for _, f := range folders {
		f := &folder{
			logger:     fs.logger,
			Name:       f.Name,
			URI:        protocol.URI(f.URI),
			controller: fs.controller,
		}
		if err := f.init(ctx); err != nil {
			return fmt.Errorf("init layer `%s`: %w", f.Name, err)
		}

		fs.folders = append(fs.folders, f)
	}

	return nil
}

type folder struct {
	logger     *slog.Logger
	Name       string
	URI        protocol.DocumentURI
	controller bufctl.Controller
	overlayfs  *file.OverlayFS
	modssource file.Source
}

func (f *folder) Souce() file.Source {
	return NewMultiSource(f.logger, f.overlayfs, f.modssource)
}

func (f *folder) Location(ctx context.Context, uri file.RelativeURI) []protocol.Location {
	f.logger.DebugContext(ctx, "location", "URI", uri)

	locs := []protocol.Location{}
	locs = append(locs, source2locations(ctx, f.overlayfs, uri)...)
	locs = append(locs, source2locations(ctx, f.modssource, uri)...)

	return locs
}

func (f *folder) init(ctx context.Context) error {
	f.logger.Debug("init folder", "Folder", f)

	ws, err := f.controller.GetWorkspace(ctx, f.URI.Filename())
	if err != nil {
		return fmt.Errorf("get workspace `%s`: %w", f.URI, err)
	}

	f.overlayfs = file.NewOverlayFS(f.URI.Filename(), file.NewDiskSource())

	mods := []file.Source{}
	for _, mod := range ws.Modules() {
		f.logger.Debug(
			"init folder module source",
			"Folder", f,
			"Module", slog.GroupValue(
				slog.Any("IsLocal", mod.IsLocal()),
				slog.String("BucketID", mod.BucketID()),
				slog.String("OpaqueID", mod.OpaqueID()),
				slog.Any("FullName", mod.ModuleFullName()),
			),
		)

		mods = append(mods, file.NewModSource(mod))
	}

	f.modssource = NewMultiSource(f.logger, mods...)
	return nil
}

func (f *folder) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("URI", f.URI),
		slog.String("Name", f.Name),
	)
}

// make sure find all related locations
func source2locations(ctx context.Context, source file.Source, uri protocol.URI) []protocol.Location {
	locs := []protocol.Location{}

	if ms, ok := source.(*MultiSource); ok {
		for _, s := range ms.Sources() {
			locs = append(locs, source2locations(ctx, s, uri)...)
		}
	} else {
		if info, err := source.Stat(ctx, uri); err == nil {
			locs = append(locs, protocol.Location{
				URI: protocol.URI(info.LocalPath()),
			})
		}
	}

	return locs
}
