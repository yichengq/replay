package replay

import (
	"fmt"
	"net/http"
)

func init() { http.HandleFunc("/share", share) }

func share(w http.ResponseWriter, r *http.Request) {

	id, err := saveSnippet(r)
	if err != nil {
		writeError(w, err)
		return
	}
	fmt.Fprint(w, id)
}
