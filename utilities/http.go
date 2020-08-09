package utilities

import (
	"net/http"
)

func StatusBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusBadRequest)
	if message == "" {
		w.Write([]byte("400 Something bad happened!"))
		return
	}
	w.Write([]byte("400 " + message))
}

func StatusOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 All reviews have been evaluated succesfully"))
}

func StatusInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Something bad happened!"))
}
