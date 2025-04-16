package cli

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func NewCmd(output io.Writer ) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "todo",
		Short: "A CLI todo app",
	}
	cmd.PersistentFlags().String("storage","mem","Choose how the todos will be stored")
	cmd.AddCommand(
		addCmd(output), 
		listCmd(output),
		getCmd(output), 
		updateCmd(output),
		deleteCmd(output),
		)
	return cmd
}

func addCmd(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "add [task]",
		Short: "Add a new todo",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			storageFlag,_:= cmd.Flags().GetString("storage")
			
			store,err:= defineStore(storageFlag)
			if err != nil {
				fmt.Fprintf(output, "%s", err)
				return
			}

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

func listCmd(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all added todos",
		Run: func(cmd *cobra.Command, args []string) {
			storageFlag,_:= cmd.Flags().GetString("storage")
			
			store,err:= defineStore(storageFlag)
			if err != nil {
				fmt.Fprintf(output, "%s", err)
				return
			}

			todos := store.GetTodos()
			for _, v := range todos {
				fmt.Fprintf(output, "%s\nAdded: %s\nCompleted: %t\n\n",
					v.Label, v.Deadline.Format("2006-01-02"), v.Completed)
			}
		},
	}
}

func getCmd(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get a todo by its id",
		Run: func(cmd *cobra.Command, args []string) {
			storageFlag,_:= cmd.Flags().GetString("storage")
			
			store,err:= defineStore(storageFlag)
			if err != nil {
				fmt.Fprintf(output, "%s", err)
				return
			}

			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Fprintf(output, "uuid not in correct format: %s", err)
				return
			}
			todo, err := store.GetTodoById(id)
			if err != nil {
				fmt.Fprintf(output, "todo not found: %s", err)
				return
			}

			fmt.Fprintf(output, "%s\nAdded: %s\nCompleted: %t\n\n",
				todo.Label, todo.Deadline.Format("2006-01-02"), todo.Completed)
		},
	}
}

func updateCmd(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "update [id] [key] [value]",
		Short: "Update a todo",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			storageFlag,_:= cmd.Flags().GetString("storage")
			
			store,err:= defineStore(storageFlag)
			if err != nil {
				fmt.Fprintf(output, "%s", err)
				return
			}

			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Fprintf(output, "uuid not in correct format: %s", err)
				return
			}
			todo, err := store.GetTodoById(id)
			if err != nil {
				fmt.Fprintf(output, "todo not found: %s", err)
				return
			}
			switch args[1] {
			case "Label":
				todo.Label = args[2]
			case "Deadline":
				todo.Deadline, err = time.Parse("2006-01-01", args[2])
				if err != nil {
					fmt.Fprintf(output, "time not in correct format: %s", err)
					return
				}
			case "Completed":
				todo.Completed, err = strconv.ParseBool(args[2])
				if err != nil {
					fmt.Fprintf(output, "completed field must be true or false: %s", err)
					return
				}
			default:
				fmt.Fprintf(output, "You must update a valid field")
			}
			store.UpdateTodo(id, todo.Label, todo.Deadline, todo.Completed)
			fmt.Fprintf(output, "Todo updated:\n%s\nAdded: %s\nCompleted: %t\n\n",
				todo.Label, todo.Deadline.Format("2006-01-02"), todo.Completed)
		},
	}
}

func deleteCmd(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a todo by its id",
		Run: func(cmd *cobra.Command, args []string) {
			storageFlag,_:= cmd.Flags().GetString("storage")
			
			store,err:= defineStore(storageFlag)
			if err != nil {
				fmt.Fprintf(output, "%s", err)
				return
			}

			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Fprintf(output, "uuid not in correct format: %s", err)
				return
			}
			err = store.DeleteTodo(id)
			if err != nil {
				fmt.Fprintf(output, "todo not found: %s", err)
				return
			}

			fmt.Fprintf(output, "Todo deleted")
		},
	}
}

func defineStore(storage string) (store.Store,error) {
	switch storage {
	case "mem":
		return store.NewInMemoryTodoStore(), nil
	case "json":
		return &store.JSONStore{},nil
	default:
		var store store.Store
		return store,errors.New("storage options: mem,json")
	}
}