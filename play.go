package replay

import (
	"html/template"
	"net/http"
)

var playTemplate = template.Must(template.ParseFiles("play.html"))

func play(w http.ResponseWriter, r *http.Request) {
	snip := getSnippet(r)
	playTemplate.Execute(w, snip)
}

func getSnippet(r *http.Request) *snippet {
	return &snippet{Body: []byte(hello)}
}

const hello = `let myVar = "Helllo";`
