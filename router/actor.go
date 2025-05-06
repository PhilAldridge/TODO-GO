package router

import (
	"time"

	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/google/uuid"
)

type StoreCommand interface {
	Execute(store store.Store)
}

type idErrReply struct {
	id  uuid.UUID
	err error
}

type todoErrReply struct {
	todo models.Todo
	err  error
}

type AddTodoCmd struct {
	StoreCommand
	label    string
	deadline time.Time
	username string
	replyCh  chan idErrReply
}

func (cmd AddTodoCmd) Execute(store store.Store) {
	id, err := store.AddTodo(cmd.label, cmd.deadline, cmd.username)
	cmd.replyCh <- idErrReply{id: id, err: err}
}

type GetTodoCmd struct {
	StoreCommand
	username string
	replyCh  chan []models.Todo
}

func (cmd GetTodoCmd) Execute(store store.Store) {
	cmd.replyCh <- store.GetTodos(cmd.username)
}

type GetTodoByIdCmd struct {
	StoreCommand
	id uuid.UUID
	username string
	replyCh  chan todoErrReply
}

func (cmd GetTodoByIdCmd) Execute(store store.Store) {
	todo,err:= store.GetTodoById(cmd.id, cmd.username)
	cmd.replyCh <- todoErrReply{todo: todo, err: err}
}

type UpdateTodoCmd struct {
	StoreCommand
	id       uuid.UUID
	field    string
	value    string
	username string
	replyCh  chan todoErrReply
}

func (cmd UpdateTodoCmd) Execute(store store.Store) {
	todo, err := store.UpdateTodo(cmd.id, cmd.field, cmd.value, cmd.username)
	cmd.replyCh <- todoErrReply{todo: todo, err: err}
}

type DeleteTodoCmd struct {
	StoreCommand
	id       uuid.UUID
	username string
	replyCh chan error
}

func (cmd DeleteTodoCmd) Execute(store store.Store) {
	cmd.replyCh <- store.DeleteTodo(cmd.id,cmd.username)
}

func StartStoreActor(store store.Store) chan<- StoreCommand {
	cmdCh := make(chan StoreCommand)

	go func() {
		for cmd := range cmdCh {
			cmd.Execute(store)
		}
	}()

	return cmdCh
}
