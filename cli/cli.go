package cli

import (
	"fmt"
	"io"
	"os/user"
	"time"

	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/spf13/cobra"
)

func getCurrentUsername() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func NewCmd(output io.Writer, store *store.TodoList) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "todo",
		Short: "A CLI todo app",
	}
	cmd.AddCommand(addCmd(output, store), listCmd(output, store))
	return cmd
}

func addCmd(output io.Writer, store *store.TodoList) *cobra.Command {
	return &cobra.Command{
		Use:   "add [task]",
		Short: "Add a new todo",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			user, err := getCurrentUsername()
			if err != nil {
				fmt.Fprintf(output, "username not found: %s", err)
				return
			}

			tm, err := time.Parse("2006-01-01", args[1])
			if err != nil {
				fmt.Fprintf(output, "time not in correct format: %s", err)
				return
			}

			uuid, err := store.AddTodo(args[0], user, tm)
			if err != nil {
				fmt.Fprintf(output, "could not add todo: %s", err)
				return
			}

			fmt.Fprintf(output, "Added: %d\n", uuid)
		},
	}
}

func listCmd(output io.Writer, store *store.TodoList) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all added todos",
		Run: func(cmd *cobra.Command, args []string) {
			todos := store.ListTodos()
			for _, v := range todos {
				fmt.Fprintf(output, "%s: %s\nAdded: %s\n\n", v.Author, v.Label, v.Deadline.Format("2006-02-01"))
			}
		},
	}
}
