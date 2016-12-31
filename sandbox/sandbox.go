package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/compile", compileHandler)
	http.HandleFunc("/_ah/health", healthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading code: %v", err), http.StatusBadRequest)
		return
	}
	output, err := compile(code, compileToJS)
	if err != nil {
		// TODO: separate internal server error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(output)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := healthCheck(); err != nil {
		http.Error(w, "Health check failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "ok")
}
