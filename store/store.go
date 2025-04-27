package store

import (
	"time"

	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
)

type Store interface {
	GetTodos(username string) []models.Todo
	GetTodoById(id uuid.UUID, username string) (models.Todo, error)
	AddTodo(
		label string,
		deadline time.Time,
		username string,
	) (uuid.UUID, error)
	UpdateTodo(
		id uuid.UUID,
		field string,
		updatedValue string,
		username string,
	) (models.Todo, error)
	DeleteTodo(id uuid.UUID, username string) error
}
