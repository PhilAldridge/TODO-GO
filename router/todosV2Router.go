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

func NewV2ApiHandler(store store.Store) TodoApiHandlerV2 {
	return TodoApiHandlerV2{actor:StartStoreActor(store)}
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
	replyCh := make(chan idErrReply)
	h.actor <- AddTodoCmd{
		label:    body.Label,
		deadline: deadline,
		username: h.username,
		replyCh:  replyCh,
	}
	reply := <-replyCh
	if reply.err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, reply.err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(reply.id.String()))
}

func (h *TodoApiHandlerV2) HandleGet(w http.ResponseWriter, r *http.Request) {
	h.setUsername(r)
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if id == "" || err != nil {
		replyCh:= make(chan []models.Todo)
		h.actor <- GetTodoCmd{
			username: h.username,
			replyCh: replyCh,
		}
		todos:=<-replyCh
		marshalAndWrite(w, todos)
		return
	}
	replyCh:= make(chan todoErrReply)
	h.actor <- GetTodoByIdCmd{
		id: uuid,
		username: h.username,
		replyCh: replyCh,
	}
	reply:=<-replyCh
	if reply.err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, reply.err), http.StatusNotFound)
		return
	}
	marshalAndWrite(w, reply.todo)
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

	replyCh:= make(chan todoErrReply)
	h.actor <- UpdateTodoCmd{
		id: uuid,
		field: body.Field,
		value: body.Value,
		username: h.username,
		replyCh: replyCh,
	}
	reply:= <-replyCh
	if reply.err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, reply.err), http.StatusNotFound)
	}
	marshalAndWrite(w, reply.todo)
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
	h.actor<- DeleteTodoCmd{
		id: uuid,
		username: h.username,
		replyCh: replyCh,
	}
	err= <- replyCh
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusNotFound)
		return
	}
	w.Write([]byte("Todo Deleted Successfully"))
}

func (h *TodoApiHandlerV2) setUsername(r *http.Request) {
	fmt.Println("and")
	
	username, ok := r.Context().Value(models.ContextKey("username")).(string)
	if !ok {
		h.username =""
	} else {
		h.username = username
	}
	fmt.Println(h.username)
}

func (h *TodoApiHandlerV2) HandleList(w http.ResponseWriter, r *http.Request) {
	h.setUsername(r)
	replyCh := make (chan []models.Todo)
	h.actor<-GetTodoCmd{
		username: h.username,
		replyCh: replyCh,
	}
	todos:=<-replyCh
	ServeTemplate("./webTemplates/list.html",todos,w)
}