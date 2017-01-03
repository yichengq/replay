package replay

import (
	"io"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

const formatURL = "https://sandbox-dot-replay-154206.appspot-preview.com/format"

func init() { http.HandleFunc("/format", format) }

func format(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Post(formatURL, "text/plain", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
