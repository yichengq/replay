package replay

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"appengine"
	"appengine/datastore"
)

const hello = `let myVar = "Hello";
print_endline myVar;`

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

func getSnippet(r *http.Request) *snippet {
	snip := &snippet{Body: []byte(hello)}
	if !strings.HasPrefix(r.URL.Path, "/p/") {
		return snip
	}
	c := appengine.NewContext(r)
	id := r.URL.Path[3:]
	key := datastore.NewKey(c, "Snippet", id, 0, nil)
	err := datastore.Get(c, key, snip)
	if err != nil {
		if err != datastore.ErrNoSuchEntity {
			c.Errorf("loading Snippet: %v", err)
		}
		return nil
	}
	return snip
}

func saveSnippet(r *http.Request) (string, error) {
	c := appengine.NewContext(r)

	var body bytes.Buffer
	_, err := io.Copy(&body, io.LimitReader(r.Body, maxSnippetSize+1))
	r.Body.Close()
	if err != nil {
		c.Errorf("reading Body: %v", err)
		return "", httpError{"Server Error", http.StatusInternalServerError}
	}
	if body.Len() > maxSnippetSize {
		return "", httpError{"Snippet is too large", http.StatusRequestEntityTooLarge}
	}

	snip := &snippet{Body: body.Bytes()}
	id := snip.ID()
	key := datastore.NewKey(c, "Snippet", id, 0, nil)
	_, err = datastore.Put(c, key, snip)
	if err != nil {
		c.Errorf("putting Snippet: %v", err)
		return "", httpError{"Server Error", http.StatusInternalServerError}
	}
	return id, nil
}
