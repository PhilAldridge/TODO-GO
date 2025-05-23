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

func (t *JSONStore) AddTodo(label string, deadline time.Time, username string) (uuid.UUID, error) {
	todos := lib.ReadJsonStore()
	store := LoadInMemoryTodoStore(todos)
	newUuid, err := store.AddTodo(label, deadline, username)
	if err == nil {
		lib.WriteJsonStore(store.GetAllTodos())
	}
	return newUuid, err
}

func (t *JSONStore) GetTodos(username string) []models.Todo {
	todos := lib.ReadJsonStore()
	store := LoadInMemoryTodoStore(todos)
	return store.GetTodos(username)
}

func (t *JSONStore) UpdateTodo(id uuid.UUID, field string, updatedValue string, username string) (models.Todo, error) {
	todos := lib.ReadJsonStore()
	store := LoadInMemoryTodoStore(todos)
	todo, err := store.UpdateTodo(id, field, updatedValue,username)
	if err == nil {
		lib.WriteJsonStore(store.GetAllTodos())
	}
	return todo, err
}

func (t *JSONStore) GetTodoById(id uuid.UUID, username string) (models.Todo, error) {
	todos := lib.ReadJsonStore()
	store := LoadInMemoryTodoStore(todos)
	return store.GetTodoById(id,username)
}

func (t *JSONStore) DeleteTodo(id uuid.UUID, username string) error {
	todos := lib.ReadJsonStore()
	store := LoadInMemoryTodoStore(todos)
	err := store.DeleteTodo(id,username)
	if err == nil {
		lib.WriteJsonStore(store.GetAllTodos())
	}
	return err
}

func(t *JSONStore) Close() {}