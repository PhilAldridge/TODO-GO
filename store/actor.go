package store

import (
	"time"

	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
)

type StoreCommand interface {
	Execute(store Store)
}

type IdErrReply struct {
	Id  uuid.UUID
	Err error
}

type TodoErrReply struct {
	Todo models.Todo
	Err  error
}

type AddTodoCmd struct {
	StoreCommand
	Label    string
	Deadline time.Time
	Username string
	ReplyCh  chan IdErrReply
}

func (cmd AddTodoCmd) Execute(store Store) {
	id, err := store.AddTodo(cmd.Label, cmd.Deadline, cmd.Username)
	cmd.ReplyCh <- IdErrReply{Id: id, Err: err}
}

type GetTodoCmd struct {
	StoreCommand
	Username string
	ReplyCh  chan []models.Todo
}

func (cmd GetTodoCmd) Execute(store Store) {
	cmd.ReplyCh <- store.GetTodos(cmd.Username)
}

type GetTodoByIdCmd struct {
	StoreCommand
	Id uuid.UUID
	Username string
	ReplyCh  chan TodoErrReply
}

func (cmd GetTodoByIdCmd) Execute(store Store) {
	todo,err:= store.GetTodoById(cmd.Id, cmd.Username)
	cmd.ReplyCh <- TodoErrReply{Todo: todo, Err: err}
}

type UpdateTodoCmd struct {
	StoreCommand
	Id       uuid.UUID
	Field    string
	Value    string
	Username string
	ReplyCh  chan TodoErrReply
}

func (cmd UpdateTodoCmd) Execute(store Store) {
	todo, err := store.UpdateTodo(cmd.Id, cmd.Field, cmd.Value, cmd.Username)
	cmd.ReplyCh <- TodoErrReply{Todo: todo, Err: err}
}

type DeleteTodoCmd struct {
	StoreCommand
	Id       uuid.UUID
	Username string
	ReplyCh chan error
}

func (cmd DeleteTodoCmd) Execute(store Store) {
	cmd.ReplyCh <- store.DeleteTodo(cmd.Id,cmd.Username)
}

type StoreActor struct {
	CommandChan chan StoreCommand
	quitChan chan struct{}
	store Store
}

func StartStoreActor(store Store) *StoreActor {
	a := &StoreActor{
		CommandChan: make(chan StoreCommand),
		quitChan: make(chan struct{}),
		store: store,
	}	
	go a.run()

	return a
}

func (a *StoreActor) run() {
	for {
		select {
		case cmd := <-a.CommandChan:
			cmd.Execute(a.store)
		case <-a.quitChan:
			return
		}
	}
}

func (a *StoreActor) Stop() {
	close(a.quitChan)
}
