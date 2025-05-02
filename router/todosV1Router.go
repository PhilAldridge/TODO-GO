package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/google/uuid"
)

func NewV1ApiHandler(store store.Store) TodoApiHandler {
	return TodoApiHandler{store: store}
}

func (h *TodoApiHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
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
	todoId, err := h.store.AddTodo(body.Label, deadline, "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(todoId.String()))
}

func (h *TodoApiHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		todos := h.store.GetTodos("")
		marshalAndWrite(w, todos)
		return
	}
	todo, err := h.store.GetTodoById(uuid,"")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	marshalAndWrite(w, todo)
}

func (h *TodoApiHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
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

	todo, err := h.store.UpdateTodo(uuid, body.Field, body.Value,"")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
	}
	marshalAndWrite(w, todo)
}

func (h *TodoApiHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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

	err = h.store.DeleteTodo(uuid,"")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	w.Write([]byte("Todo Deleted Successfully"))
}

func (h *TodoApiHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	todos := h.store.GetTodos("")
	ServeTemplate("./webTemplates/list.html",todos,w)
}

