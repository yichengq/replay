package replay

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/memcache"
	"appengine/urlfetch"
)

const compileURL = "https://sandbox-dot-replay-154206.appspot-preview.com/compile?type="

func init() { http.HandleFunc("/compile", compile) }

type result struct {
	Status  int
	Message []byte
}

func compile(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	compileType := r.URL.Query().Get("type")
	if compileType == "" {
		compileType = "to_run"
	}

	body := io.Reader(r.Body)
	var input []byte
	if r.ContentLength >= 0 && r.ContentLength < 250 && compileType == "to_run" {
		var err error
		input, err = ioutil.ReadAll(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// check cached result
		if item, err := memcache.Get(ctx, string(input)); err == nil {
			var r result
			if err := json.Unmarshal(item.Value, &r); err == nil {
				w.WriteHeader(r.Status)
				w.Write(r.Message)
				return
			}
		}

		body = bytes.NewBuffer(input)
	}

	resp, err := client.Post(compileURL+compileType, "text/plain", body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.ContentLength >= 0 && resp.ContentLength < 500000 {
		output, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if bs, err := json.Marshal(&result{
			Status:  resp.StatusCode,
			Message: output,
		}); err == nil {
			memcache.Set(ctx, &memcache.Item{
				Key:   string(input),
				Value: bs,
			})
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(output)
		return
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
