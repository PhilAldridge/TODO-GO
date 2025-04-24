package store

import (
	"errors"
	"strconv"
	"time"

	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
)

type InMemoryStore struct {
	Store
	todos []models.Todo
}

func NewInMemoryTodoStore() *InMemoryStore {
	s := InMemoryStore{}
	return &s
}

func LoadInMemoryTodoStore(todos []models.Todo) *InMemoryStore {
	s := InMemoryStore{todos:todos}
	return &s
}

func (t *InMemoryStore) AddTodo(label string, deadline time.Time) (uuid.UUID, error) {
	newUuid := uuid.New()

	t.todos = append(t.todos, models.Todo{
		Id:        newUuid,
		Label:     label,
		Completed: false,
		Deadline:  deadline,
	})
	return newUuid, nil
}

func (t *InMemoryStore) GetTodos() []models.Todo {
	return t.todos
}

func (t *InMemoryStore) UpdateTodo(id uuid.UUID, field string, updatedValue string) (models.Todo, error) {
	for i, todo := range t.todos {
		if todo.Id == id {
			switch field {
			case "label":
				t.todos[i].Label = updatedValue
			case "deadline":
				newDeadline, err:= time.Parse("2006-01-02",updatedValue)
				if err != nil {
					return models.Todo{},err
				}
				t.todos[i].Deadline = newDeadline
			case "completed":
				newCompleted, err:= strconv.ParseBool(updatedValue)
				if err != nil {
					return models.Todo{},err
				}
				t.todos[i].Completed = newCompleted
			default:
				return models.Todo{}, errors.New("allowed update fields: label, deadline, completed")
			}
			return t.todos[i], nil
		}
	}
	return models.Todo{}, errors.New("todo not found")
}

func (t *InMemoryStore) GetTodoById(id uuid.UUID) (models.Todo, error) {
	for _, todo := range t.todos {
		if todo.Id == id {
			return todo, nil
		}
	}
	return models.Todo{}, errors.New("todo not found")
}

func (t *InMemoryStore) DeleteTodo(id uuid.UUID) error {
	for i, todo := range t.todos {
		if todo.Id == id {
			t.todos = append(t.todos[:i], t.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found")
}
