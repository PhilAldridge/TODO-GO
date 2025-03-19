package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/spf13/cobra"
)

func NewCmd(output io.Writer, store store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "todo",
		Short: "A CLI todo app",
	}
	cmd.AddCommand(addCmd(output, store), listCmd(output, store))
	return cmd
}

func addCmd(output io.Writer, store store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "add [task]",
		Short: "Add a new todo",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			tm, err := time.Parse("2006-01-01", args[1])
			if err != nil {
				fmt.Fprintf(output, "time not in correct format: %s", err)
				return
			}

			uuid, err := store.AddTodo(args[0], tm)
			if err != nil {
				fmt.Fprintf(output, "could not add todo: %s", err)
				return
			}

			fmt.Fprintf(output, "Added: %d\n", uuid)
		},
	}
}

func listCmd(output io.Writer, store store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all added todos",
		Run: func(cmd *cobra.Command, args []string) {
			todos := store.GetTodos()
			for _, v := range todos {
				fmt.Fprintf(output, "%s\nAdded: %s\nCompleted: %t\n\n",
					v.Label, v.Deadline.Format("2006-01-02"), v.Completed)
			}
		},
	}
}
