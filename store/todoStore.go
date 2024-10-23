package store

import (
	"fmt"
	"io"
	"time"
	"github.com/google/uuid"
)

type TodoList struct {
	todos []Todo
}

type Todo struct {
	id        string
	label     string
	author    string
	completed bool
	deadline  time.Time
}

func NewTodoList() *TodoList {
	t := TodoList{
		todos: []Todo{},
	}
	return &t
}

func (t *TodoList) AddTodo(label string, author string, deadline time.Time) {
	t.todos = append(t.todos, Todo{
		id:        uuid.NewString(),
		label:     label,
		author:    author,
		completed: false,
		deadline:  deadline,
	})
}

func (t *TodoList) ListTodos(writer io.Writer) {
	for _, todo := range t.todos {
		completionSymbol := "☐"
		if todo.completed {
			completionSymbol = "☑"
		}
		fmt.Fprintf(writer, "%s - Due by: %s - %s\n", todo.label, todo.deadline.Format("01/02/2006"), completionSymbol)
	}
}
