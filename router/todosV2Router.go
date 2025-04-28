package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/google/uuid"
)

func NewV2ApiHandler(store store.Store) TodoApiHandlerV2 {
	return TodoApiHandlerV2{store: store}
}

func (h *TodoApiHandlerV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.username = r.Context().Value("username").(string)

	switch {
	case r.Method == http.MethodPut:
		h.handlePut(w, r)
	case r.Method == http.MethodGet:
		h.handleGet(w, r)
	case r.Method == http.MethodPatch:
		h.handlePatch(w, r)
	case r.Method == http.MethodDelete:
		h.handleDelete(w, r)
	default:
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
}

func (h *TodoApiHandlerV2) handlePut(w http.ResponseWriter, r *http.Request) {
	var body V1PutBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if body.Label == "" || err != nil {
		http.Error(w, "Put must include a todo label and a deadline (in the form 2006-01-02)", http.StatusBadRequest)
		return
	}
	deadline, err := time.Parse("2006-01-02", body.Deadline)
	if err != nil {
		http.Error(w, "Put must include a todo label and a deadline (in the form 2006-01-02)", http.StatusBadRequest)
		return
	}
	todoId, err := h.store.AddTodo(body.Label, deadline,h.username)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(todoId.String()))
}

func (h *TodoApiHandlerV2) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		todos := h.store.GetTodos(h.username)
		marshalAndWrite(w, todos)
		return
	}
	todo, err := h.store.GetTodoById(uuid, h.username)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	marshalAndWrite(w, todo)
}

func (h *TodoApiHandlerV2) handlePatch(w http.ResponseWriter, r *http.Request) {
	var body V1PatchBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Provide Id, Field, Value", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(body.Id)
	if body.Id == "" || err != nil {
		http.Error(w, "Error: patch method requires a valid id", http.StatusBadRequest)
		return
	}

	todo, err := h.store.UpdateTodo(uuid, body.Field, body.Value, h.username)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound) 
	}
	marshalAndWrite(w, todo)
}

func (h *TodoApiHandlerV2) handleDelete(w http.ResponseWriter, r *http.Request) {
	var body V1DeleteBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Provide a valid id", http.StatusBadRequest)
		return
	}
	uuid, err := uuid.Parse(body.Id)
	if body.Id == "" || err != nil {
		http.Error(w, "Error: delete method requires a valid id", http.StatusBadRequest)
		return
	}

	err = h.store.DeleteTodo(uuid, h.username)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	w.Write([]byte("Todo Deleted Successfully"))
}
