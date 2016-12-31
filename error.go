package replay

import (
	"net/http"
)

type httpError struct {
	Message    string
	StatusCode int
}

func (err httpError) Error() string { return err.Message }

func writeError(w http.ResponseWriter, err error) {
	msg := "Server Error"
	code := http.StatusInternalServerError
	if e, ok := err.(httpError); ok {
		msg = e.Message
		code = e.StatusCode
	}
	http.Error(w, msg, code)
}
