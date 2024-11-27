package file

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type RelativeURI = protocol.URI

type ErrInvalidURIReason string

const (
	ErrInvalidURIReasonOnlyAcceptRelativeURI ErrInvalidURIReason = "only accept RelativeURI"
	ErrInvalidURIReasonOnlyAcceptDocumentURI ErrInvalidURIReason = "only accept DocumentURI"
)

type ErrInvalidURI struct {
	uri    protocol.URI
	reason ErrInvalidURIReason
}

func ErrorInvalidURI(uri protocol.URI, reason ErrInvalidURIReason) error {
	return &ErrInvalidURI{
		uri:    uri,
		reason: reason,
	}
}

func (e *ErrInvalidURI) Error() string {
	return fmt.Sprintf("invalid URI: `%s`, %s", e.uri, e.reason)
}

func URI2Relative(u protocol.URI, base string) (RelativeURI, error) {
	filename := URI2Filename(u)

	// just relative uri
	// TODO: check if match base?
	if !filepath.IsAbs(filename) {
		return RelativeURI(filename), nil
	}

	rel, err := filepath.Rel(base, filename)
	if err != nil {
		return "", err
	}

	return uri.New(rel), nil
}

func URI2Filename(u protocol.URI) string {
	if strings.HasPrefix(string(u), uri.FileScheme) {
		return u.Filename()
	} else {
		// as relative or just path

		return string(u)
	}
}

// file schemed for abs path
// stay the same for relative path
func NormalURI(u protocol.URI) protocol.URI {
	return NormalURIStr(string(u))
}

// file schemed for abs path
// stay the same for relative path
func NormalURIStr(u string) protocol.URI {
	if strings.HasPrefix(u, uri.FileScheme) {
		// file schemed

		return uri.New(u)
	} else if filepath.IsAbs(u) {
		// abs path

		u := url.URL{
			Scheme: uri.FileScheme,
			Path:   u,
		}

		return uri.New(u.String())
	} else {
		// relative path

		return RelativeURI(u)
	}
}

func IsDocumentURI(u protocol.URI) bool {
	return strings.HasPrefix(string(u), uri.FileScheme) && filepath.IsAbs(u.Filename())
}

func TrimScheme(u protocol.URI) string {
	return filepath.FromSlash(
		strings.ReplaceAll(
			string(u),
			uri.FileScheme+"://",
			"",
		),
	)
}
