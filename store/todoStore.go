package store

import (
	"time"

	"github.com/google/uuid"
)

type TodoList struct {
	todos map[uuid.UUID]Todo
}

type Todo struct {
	Label     string
	Author    string
	Completed bool
	Deadline  time.Time
}

func NewTodoList() *TodoList {
	t := TodoList{
		todos: make(map[uuid.UUID]Todo),
	}
	return &t
}

func (t *TodoList) AddTodo(label string, author string, deadline time.Time) {
	t.todos[uuid.New()] = Todo{
		Label:     label,
		Author:    author,
		Completed: false,
		Deadline:  deadline,
	}
}

func (t *TodoList) ListTodos() {

}

func (t *TodoList) PatchTodo(id uuid.UUID, label string, deadline time.Time, completed bool) (Todo, error) {
	t.todos[id] = Todo{
		Label:     label,
		Author:    t.todos[id].Author,
		Completed: completed,
		Deadline:  deadline,
	}
	return t.todos[id], nil
}
