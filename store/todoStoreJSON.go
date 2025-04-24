package store

import (
	"time"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
)

type JSONStore struct {
	Store
}


func (t *JSONStore) AddTodo(label string, deadline time.Time) (uuid.UUID, error) {
	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	newUuid,err := store.AddTodo(label,deadline)
	if err== nil {
		lib.WriteJson(store.GetTodos())
	}
	return newUuid, err
}

func (t *JSONStore) GetTodos() []models.Todo {
	todos := lib.ReadJson()

	return todos
}

func (t *JSONStore) UpdateTodo(id uuid.UUID, field string, updatedValue string) (models.Todo, error) {
	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	todo,err:= store.UpdateTodo(id,field,updatedValue)
	if err== nil {
		lib.WriteJson(store.GetTodos())
	}
	return todo,err
}

func (t *JSONStore) GetTodoById(id uuid.UUID) (models.Todo, error) {
	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	return store.GetTodoById(id)
}

func (t *JSONStore) DeleteTodo(id uuid.UUID) error {
	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	err:= store.DeleteTodo(id)
	if err== nil {
		lib.WriteJson(store.GetTodos())
	}
	return err
}
