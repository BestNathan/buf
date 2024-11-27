package file

import "go.lsp.dev/protocol"

type simplefile struct {
	uri     protocol.DocumentURI
	content []byte
	err     error
}

func (s *simplefile) URI() protocol.DocumentURI {
	return s.uri
}

func (s *simplefile) Version() int32 {
	return -1
}

func (s *simplefile) Content() ([]byte, error) {
	return s.content, s.err
}
