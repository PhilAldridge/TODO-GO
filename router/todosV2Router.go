package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/google/uuid"
)

func NewV2ApiHandler(actor *store.StoreActor) TodoApiHandlerV2 {
	return TodoApiHandlerV2{actor:actor}
}

func (h *TodoApiHandlerV2) HandlePut(w http.ResponseWriter, r *http.Request) {
	h.setUsername(r)
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
	replyCh := make(chan store.IdErrReply)
	h.actor.CommandChan <- store.AddTodoCmd{
		Label:    body.Label,
		Deadline: deadline,
		Username: h.username,
		ReplyCh:  replyCh,
	}
	reply := <-replyCh
	if reply.Err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, reply.Err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(reply.Id.String()))
}

func (h *TodoApiHandlerV2) HandleGet(w http.ResponseWriter, r *http.Request) {
	h.setUsername(r)
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		replyCh:= make(chan []models.Todo)
		h.actor.CommandChan <- store.GetTodoCmd{
			Username: h.username,
			ReplyCh: replyCh,
		}
		todos:=<-replyCh
		marshalAndWrite(w, todos)
		return
	}
	replyCh:= make(chan store.TodoErrReply)
	h.actor.CommandChan <- store.GetTodoByIdCmd{
		Id: uuid,
		Username: h.username,
		ReplyCh: replyCh,
	}
	reply:=<-replyCh
	if reply.Err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, reply.Err), http.StatusNotFound)
		return
	}
	marshalAndWrite(w, reply.Todo)
}

func (h *TodoApiHandlerV2) HandlePatch(w http.ResponseWriter, r *http.Request) {
	h.setUsername(r)
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

	replyCh:= make(chan store.TodoErrReply)
	h.actor.CommandChan <- store.UpdateTodoCmd{
		Id: uuid,
		Field: body.Field,
		Value: body.Value,
		Username: h.username,
		ReplyCh: replyCh,
	}
	reply:= <-replyCh
	if reply.Err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, reply.Err), http.StatusNotFound)
	}
	marshalAndWrite(w, reply.Todo)
}

func (h *TodoApiHandlerV2) HandleDelete(w http.ResponseWriter, r *http.Request) {
	h.setUsername(r)
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

	replyCh:= make(chan error)
	h.actor.CommandChan <- store.DeleteTodoCmd{
		Id: uuid,
		Username: h.username,
		ReplyCh: replyCh,
	}
	err= <- replyCh
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	w.Write([]byte("Todo Deleted Successfully"))
}

func (h *TodoApiHandlerV2) setUsername(r *http.Request) {
	username, ok := r.Context().Value(models.ContextKey("username")).(string)
	if !ok {
		h.username =""
	} else {
		h.username = username
	}
}

func (h *TodoApiHandlerV2) HandleList(w http.ResponseWriter, r *http.Request) {
	h.setUsername(r)
	replyCh := make (chan []models.Todo)
	h.actor.CommandChan<-store.GetTodoCmd{
		Username: h.username,
		ReplyCh: replyCh,
	}
	todos:=<-replyCh
	ServeTemplate("./webTemplates/list.html",todos,w)
}