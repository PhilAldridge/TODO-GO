package store

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id		  uuid.UUID
	Label     string
	Author    string
	Completed bool
	Deadline  time.Time
}

type Store interface {
	GetTodos() []Todo
	GetTodoById(id uuid.UUID) (Todo,error)
	AddTodo(
		label string,
		deadline time.Time,
	) (uuid.UUID, error)
	UpdateTodo(
		id uuid.UUID,
		label string,
		deadline time.Time,
		completed bool,
	) (Todo,error)
	DeleteTodo(id uuid.UUID) error
}
