package router

import (
	"fmt"
	"net/http"
	"time"

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
		h.handlePut(w,r)
	case r.Method == http.MethodGet:
		h.handleGet(w,r)
	case r.Method == http.MethodPatch:
		h.handlePatch(w,r)
	case r.Method == http.MethodDelete:
		h.handleDelete(w,r)
	default:
		http.Error(w,"Invalid request", http.StatusBadRequest)
	}
}

func (h *ApiHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	label:= r.Form.Get("label")
	deadlineString:= r.Form.Get("deadline")
	deadline,err:= time.Parse("2006-01-02",deadlineString)
	if label =="" || err!=nil {
		http.Error(w,"Put must include a todo label and a deadline (in the form 2006-01-02)",http.StatusBadRequest)
		return
	}
	todoId, err:= h.store.AddTodo(label,deadline)
	if err != nil {
		http.Error(w,fmt.Sprintf(`{"error": "%s"}`, err),2)//TODO
		return
	}
	//marshalAndWrite(w,todoId)
	w.Write([]byte(todoId.String()))
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
		http.Error(w,fmt.Sprintf(`{"error": "%s"}`, err),2)//TODO
		return
	}
	marshalAndWrite(w,todo)
}

func (h *ApiHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id:= r.Form.Get("id")
	field:= r.Form.Get("field")
	updatedValue:= r.Form.Get("value")
	
	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		http.Error(w,"Error: patch method requires a valid id",http.StatusBadRequest)
		return
	}

	todo,err:= h.store.UpdateTodo(uuid,field, updatedValue)
	if err!=nil {
		http.Error(w,fmt.Sprintf(`{"error": "%s"}`, err),http.StatusNotFound)//TODO
	}
	marshalAndWrite(w,todo)
}

func (h *ApiHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id:= r.Form.Get("id")

	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		http.Error(w,"Error: delete method requires a valid id",http.StatusBadRequest)
		return
	}

	err = h.store.DeleteTodo(uuid)
	if err != nil {
		http.Error(w,fmt.Sprintf(`{"error": "%s"}`, err),http.StatusNotFound)//TODO
		return
	}
	w.Write([]byte("Todo Deleted Successfully"))
}
