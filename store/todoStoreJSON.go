package store

import (
	"errors"
	"strconv"
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

	newUuid := uuid.New()

	todos = append(todos, models.Todo{
		Id: newUuid,
		Label:     label,
		Completed: false,
		Deadline:  deadline,
	})

	lib.WriteJson(todos)
	return newUuid, nil
}

func (t *JSONStore) GetTodos() []models.Todo {
	todos := lib.ReadJson()

	return todos
}

func (t *JSONStore) UpdateTodo(id uuid.UUID, field string, updatedValue string) (models.Todo, error) {
	todos := lib.ReadJson()

	for i,todo:= range todos {
		if todo.Id == id {
			switch field {
			case "label":
				todos[i].Label = updatedValue
			case "deadline":
				newDeadline, err:= time.Parse("2006-01-01",updatedValue)
				if err != nil {
					return models.Todo{},err
				}
				todos[i].Deadline = newDeadline
			case "completed":
				newCompleted, err:= strconv.ParseBool(updatedValue)
				if err != nil {
					return models.Todo{},err
				}
				todos[i].Completed = newCompleted
			default:
				return models.Todo{}, errors.New("allowed update fields: label, deadline, completed")
			}
			lib.WriteJson(todos)
			return todos[i],nil
		}
	}
	return models.Todo{}, errors.New("todo not found")
}

func (t *JSONStore) GetTodoById(id uuid.UUID) (models.Todo, error) {
	todos := lib.ReadJson()

	for _,todo:=range todos {
		if todo.Id == id {
			return todo,nil
		}
	}
	return models.Todo{}, errors.New("todo not found")
}

func (t *JSONStore) DeleteTodo(id uuid.UUID) error {
	todos := lib.ReadJson()
	
	for i,todo:= range todos {
		if todo.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			lib.WriteJson(todos)
			return nil
		}
	}
	return errors.New("todo not found")
}
