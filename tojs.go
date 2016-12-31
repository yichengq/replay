package replay

import (
	"io/ioutil"
	"net/http"
)

func toJS(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "XXX", http.StatusBadRequest)
		return
	}
	out := compile(d, compileToJS)
	w.Write(out)
	return
}
