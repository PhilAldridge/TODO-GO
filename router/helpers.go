package router

import (
	"encoding/json"
	"net/http"
)

func internalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func writeJSONResponse(w http.ResponseWriter,statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func marshalAndWrite[T any](w http.ResponseWriter, data T) {
	resp, err := json.Marshal(data)
	if err != nil {
		internalServerErrorHandler(w)
		return
	}
	writeJSONResponse(w, http.StatusOK,resp)
}