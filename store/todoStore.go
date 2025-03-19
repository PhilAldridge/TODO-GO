package store

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type InMemoryStore struct {
	Store
	todos []Todo
}

func NewTodoStore() *InMemoryStore {
	s:= InMemoryStore{}
	return &s
}

func (t *InMemoryStore) AddTodo(label string, deadline time.Time) (uuid.UUID, error) {
	newUuid := uuid.New()

	t.todos = append(t.todos, Todo{
		Id: newUuid,
		Label:     label,
		Completed: false,
		Deadline:  deadline,
	})
	return newUuid, nil
}

func (t *InMemoryStore) GetTodos() []Todo {
	return t.todos
}

func (t *InMemoryStore) UpdateTodo(id uuid.UUID, label string, deadline time.Time, completed bool) (Todo, error) {
	for i,todo:= range t.todos {
		if todo.Id == id {
			t.todos[i] = Todo{
				Id: id,
				Label: label,
				Deadline: deadline,
				Completed: completed,
			}
			return t.todos[i],nil
		}
	}
	return Todo{}, errors.New("todo not found")
}

func (t *InMemoryStore) GetTodoById(id uuid.UUID) (Todo, error) {
	for _,todo:=range t.todos {
		if todo.Id == id {
			return todo,nil
		}
	}
	return Todo{}, errors.New("todo not found")
}

func (t *InMemoryStore) DeleteTodo(id uuid.UUID) error {
	for i,todo:= range t.todos {
		if todo.Id == id {
			t.todos = append(t.todos[:i], t.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found")
}
