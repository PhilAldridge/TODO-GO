package router

import (
	"fmt"
	"net/http"

	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/google/uuid"
)

type ApiHandler struct{
	store store.Store
}

func NewApiHandler(store store.Store) ApiHandler {
	return ApiHandler{store: store}
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPut:
		h.CreateTodo(w,r)
	case r.Method == http.MethodGet:
		h.handleGet(w,r)
	case r.Method == http.MethodPatch:
		h.UpdateTodo(w,r)
	case r.Method == http.MethodDelete:
		h.DeleteTodo(w,r)
	default:
		http.Error(w,"Invalid request", http.StatusBadRequest)
	}
}

func (h *ApiHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	
}

func (h *ApiHandler) handleGet(w http.ResponseWriter, r *http.Request)    {
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		todos:= h.store.GetTodos()
		marshalAndWrite(w,todos)
		return
	}
	todo, err:= h.store.GetTodoById(uuid)
	if err != nil {
		writeJSONResponse(w, 3, []byte(fmt.Sprintf(`{"error": "%s"}`, err))) //TODO
		return
	}
	marshalAndWrite(w,todo)
}

func (h *ApiHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {}

func (h *ApiHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {}
