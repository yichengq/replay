package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/format", formatHandler)
	http.HandleFunc("/compile", compileHandler)
	http.HandleFunc("/_ah/health", healthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type result struct {
	output []byte
	errStr string
}

func formatHandler(w http.ResponseWriter, r *http.Request) {
	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading code: %v", err), http.StatusBadRequest)
		return
	}
	res, err := formatReason(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.errStr != "" {
		http.Error(w, res.errStr, http.StatusBadRequest)
		return
	}
	w.Write(res.output)
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading code: %v", err), http.StatusBadRequest)
		return
	}
	ct, err := parseCompileType(r.FormValue("type"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid compile type %q: %v", r.FormValue("type"), err), http.StatusBadRequest)
		return
	}
	res, err := compile(code, ct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.errStr != "" {
		http.Error(w, res.errStr, http.StatusBadRequest)
		return
	}
	w.Write(res.output)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := healthCheck(); err != nil {
		http.Error(w, "Health check failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "ok")
}
