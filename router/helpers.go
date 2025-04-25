package router

import (
	"encoding/json"
	"net/http"

	"github.com/PhilAldridge/TODO-GO/store"
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

type V1PutBody struct {
	Label    string
	Deadline string
}

type V1PatchBody struct {
	Id    string
	Field string
	Value string
}

type V1DeleteBody struct {
	Id string
}

type TodoApiHandler struct {
	store store.Store
}

type UserPutBody struct {
	Username string
	Password string
}
