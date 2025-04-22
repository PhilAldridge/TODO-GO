package router

import (
	"encoding/json"
	"net/http"
)

func internalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func marshalAndWrite[T any](w http.ResponseWriter, data T) {
	resp, err := json.Marshal(data)
	if err != nil {
		internalServerErrorHandler(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}