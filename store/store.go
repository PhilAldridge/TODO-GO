package store

import (
	"time"

	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/google/uuid"
)

type Store interface {
	GetTodos() []models.Todo
	GetTodoById(id uuid.UUID) (models.Todo, error)
	AddTodo(
		label string,
		deadline time.Time,
	) (uuid.UUID, error)
	UpdateTodo(
		id uuid.UUID,
		field string,
		updatedValue string,
	) (models.Todo, error)
	DeleteTodo(id uuid.UUID) error
}
