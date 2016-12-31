package replay

import (
	"html/template"
	"net/http"
)

func init() { http.HandleFunc("/", play) }

var playTemplate = template.Must(template.ParseFiles("play.html"))

func play(w http.ResponseWriter, r *http.Request) {
	snip := getSnippet(r)
	if snip == nil {
		http.Error(w, "Snippet not found", http.StatusNotFound)
	}
	playTemplate.Execute(w, snip)
}
