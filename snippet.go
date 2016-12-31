package replay

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
)

const salt = "[replace this with something unique]"

const maxSnippetSize = 64 * 1024

type snippet struct {
	Body []byte
}

func (s *snippet) ID() string {
	h := sha1.New()
	io.WriteString(h, salt)
	h.Write(s.Body)
	sum := h.Sum(nil)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(sum)))
	base64.URLEncoding.Encode(b, sum)
	return string(b)[:10]
}
