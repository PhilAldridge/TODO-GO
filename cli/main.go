package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/models"
	"github.com/PhilAldridge/TODO-GO/router"
	"github.com/spf13/cobra"
)

var url string

func main() {
	lib.LoadConfig("../.env")
	url = fmt.Sprintf("http://localhost%s/Todos",lib.PortNo)
	cmd := NewCmd(os.Stdout)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

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
			body,err:= json.Marshal(router.PutBody{
				Label: args[0],
				Deadline: args[1],
			})

			if err!= nil {
				fmt.Println("Error: add needs two arguments, label and deadline")
				return
			}
			sendAndReceive(http.MethodPut,body)
		},
	}
}

func listCmd(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all added todos",
		Run: func(cmd *cobra.Command, args []string) {
			res,err:= http.Get(url)
			if err!= nil {
				fmt.Printf("error making http request: %s\n", err)
  				return
			}
			var todos []models.Todo

			err = json.NewDecoder(res.Body).Decode(&todos)
			if err!=nil {
				fmt.Printf("error making http request: %s\n", err)
  				return
			}
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
			res,err:= http.Get(fmt.Sprintf("%s?id=%s",url,args[0]))
			if err!= nil {
				fmt.Printf("error making http request: %s\n", err)
  				return
			}
			var todo models.Todo

			err = json.NewDecoder(res.Body).Decode(&todo)
			if err!=nil {
				fmt.Printf("error making http request: %s\n", err)
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
			body,err:= json.Marshal(router.PatchBody{
				Id: args[0],
				Field: args[1],
				Value: args[2],
			})

			if err!= nil {
				fmt.Println("Error: update needs three arguments, id, field and value")
				return
			}

			sendAndReceive(http.MethodPatch,body)
		},
	}
}

func deleteCmd(output io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a todo by its id",
		Run: func(cmd *cobra.Command, args []string) {
			body,err:= json.Marshal(router.DeleteBody{
				Id: args[0],
			})

			if err!= nil {
				fmt.Println("Error: delete needs one argument, id")
				return
			}
			sendAndReceive(http.MethodDelete,body)
		},
	}
}

func sendAndReceive(method string, body []byte) {
	req,err:= http.NewRequest(method,url,bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return
	   }

	res, err := http.DefaultClient.Do(req)
	 if err != nil {
		  fmt.Printf("client: error making http request: %s\n", err)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return
	}
	fmt.Printf("client: response body: %s\n", resBody)
}