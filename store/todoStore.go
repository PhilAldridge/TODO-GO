package store

import (
	"errors"
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

func (t *TodoList) AddTodo(label string, author string, deadline time.Time) (uuid.UUID, error) {
	newUuid := uuid.New()

	t.todos[newUuid] = Todo{
		Label:     label,
		Author:    author,
		Completed: false,
		Deadline:  deadline,
	}
	return newUuid, nil
}

func (t *TodoList) ListTodos() map[uuid.UUID]Todo {
	return t.todos
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

func (t *TodoList) GetTodo(id uuid.UUID) (Todo, error) {
	todo, ok := t.todos[id]
	if !ok {
		return Todo{}, errors.New("todo not found")
	}
	return todo, nil
}

func (t *TodoList) DeleteTodo(id uuid.UUID) error {
	if _, ok := t.todos[id]; !ok {
		return errors.New("todo not found")
	}
	delete(t.todos, id)
	return nil
}
