package store

import (
	"sync"
	"time"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
)

type JSONStore struct {
	Store
	mutex sync.Mutex
}


func (t *JSONStore) AddTodo(label string, deadline time.Time) (uuid.UUID, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	newUuid,err := store.AddTodo(label,deadline)
	if err== nil {
		lib.WriteJson(store.GetTodos())
	}
	return newUuid, err
}

func (t *JSONStore) GetTodos() []models.Todo {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	todos := lib.ReadJson()

	return todos
}

func (t *JSONStore) UpdateTodo(id uuid.UUID, field string, updatedValue string) (models.Todo, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	todo,err:= store.UpdateTodo(id,field,updatedValue)
	if err== nil {
		lib.WriteJson(store.GetTodos())
	}
	return todo,err
}

func (t *JSONStore) GetTodoById(id uuid.UUID) (models.Todo, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	return store.GetTodoById(id)
}

func (t *JSONStore) DeleteTodo(id uuid.UUID) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	todos := lib.ReadJson()
	store:= LoadInMemoryTodoStore(todos)
	err:= store.DeleteTodo(id)
	if err== nil {
		lib.WriteJson(store.GetTodos())
	}
	return err
}
