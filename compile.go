package replay

import (
	"io"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

const compileURL = "https://sandbox-dot-replay-154206.appspot-preview.com/compile?type=to_js"

func init() { http.HandleFunc("/compile", compile) }

func compile(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Post(compileURL, "text/plain", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
