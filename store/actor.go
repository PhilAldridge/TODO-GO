package store

import (
	"time"

	"github.com/google/uuid"
)

type StoreCommand interface {
	Execute(store Store)
}

type IdErrReply struct {
	Id uuid.UUID
	Error error
}

type AddTodoCmd struct {
	StoreCommand
	Label string
	Deadline time.Time
	Username string
	ReplyCh chan IdErrReply
}

func (cmd AddTodoCmd) Execute(store Store) {
	id,err:= store.AddTodo(cmd.Label, cmd.Deadline, cmd.Username) 
	cmd.ReplyCh <- IdErrReply{Id:id, Error:err}
}

func StartStoreActor(store Store) chan <- StoreCommand {
	cmdCh:= make(chan StoreCommand)

	go func() {
		for cmd:= range cmdCh {
			cmd.Execute(store)
		}
	}()

	return cmdCh
}