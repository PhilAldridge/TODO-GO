package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/google/uuid"
)

// Really this should be a static variable in the store rather than here.
var store3 store.Store

func SetupTodo(store store.Store) {
	store3 = store
}

func HandlePut3(w http.ResponseWriter, r *http.Request) {
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
	todoId, err := store3.AddTodo(body.Label, deadline, "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(todoId.String()))
}

func HandleGet3(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		todos := store3.GetTodos("")
		marshalAndWrite(w, todos)
		return
	}
	todo, err := store3.GetTodoById(uuid, "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	marshalAndWrite(w, todo)
}

func HandlePatch3(w http.ResponseWriter, r *http.Request) {
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

	todo, err := store3.UpdateTodo(uuid, body.Field, body.Value, "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
	}
	marshalAndWrite(w, todo)
}

func HandleDelete3(w http.ResponseWriter, r *http.Request) {
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

	err = store3.DeleteTodo(uuid, "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	w.Write([]byte("Todo Deleted Successfully"))
}
