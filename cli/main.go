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

var (
	url string
	username string 
	password string 
	jwt string
	cmd *cobra.Command = NewCmd()
)

func main() {
	lib.LoadConfig("../.env")
	cobra.OnInitialize(login)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "todo",
		Short: "A CLI todo app",
	}
	cmd.PersistentFlags().StringVar(&username,"username", "", "Set user for v2 api usage")
	cmd.PersistentFlags().StringVar(&password,"password","","Set password for v2 api usage")
	cmd.AddCommand(
		addCmd(),
		listCmd(),
		getCmd(),
		updateCmd(),
		deleteCmd(),
	)
	return cmd
}

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [task]",
		Short: "Add a new todo",
		Args:  cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			body, err := json.Marshal(router.V1PutBody{
				Label:    args[0],
				Deadline: args[1],
			})

			if err != nil {
				fmt.Println("Error: add needs two arguments, label and deadline")
				return
			}
			res:= sendAndReceive(http.MethodPut, body,url)
			if len(res) >0 {
				fmt.Printf("client: response body: %s\n", res)
			}
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all added todos",
		Run: func(cmd *cobra.Command, args []string) {
			res:= sendAndReceive(http.MethodGet,[]byte{},url)

			var todos []models.Todo

			err:= json.Unmarshal(res,&todos)
			if err != nil {
				fmt.Printf("error making http request: %s\n", err)
				return
			}
			for _, v := range todos {
				fmt.Printf("%s\nAdded: %s\nCompleted: %t\n\n",
					v.Label, v.Deadline.Format("2006-01-02"), v.Completed)
			}
		},
	}
}

func getCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get a todo by its id",
		Run: func(cmd *cobra.Command, args []string) {
			res:= sendAndReceive(http.MethodGet,[]byte{},url+"?id="+args[0])
			var todo models.Todo

			err := json.Unmarshal(res,&todo)
			if err != nil {
				fmt.Printf("error making http request: %s\n", err)
				return
			}
			fmt.Printf("%s\nAdded: %s\nCompleted: %t\n\n",
				todo.Label, todo.Deadline.Format("2006-01-02"), todo.Completed)
		},
	}
}

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [id] [key] [value]",
		Short: "Update a todo",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			body, err := json.Marshal(router.V1PatchBody{
				Id:    args[0],
				Field: args[1],
				Value: args[2],
			})

			if err != nil {
				fmt.Println("Error: update needs three arguments, id, field and value")
				return
			}

			res:= sendAndReceive(http.MethodPatch, body,url)
			if len(res) >0 {
				fmt.Printf("client: response body: %s\n", res)
			}
		},
	}
}

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a todo by its id",
		Run: func(cmd *cobra.Command, args []string) {
			body, err := json.Marshal(router.V1DeleteBody{
				Id: args[0],
			})

			if err != nil {
				fmt.Println("Error: delete needs one argument, id")
				return
			}
			res:=sendAndReceive(http.MethodDelete, body,url)
			if len(res) >0 {
				fmt.Printf("client: response body: %s\n", res)
			}
		},
	}
}

func sendAndReceive(method string, body []byte, urlString string) []byte {
	req, err := http.NewRequest(method, urlString, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return []byte{}
	}

	req.Header.Set("Authorization","Bearer "+jwt)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return []byte{}
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return []byte{}
	}

	if res.StatusCode == http.StatusOK {
		return resBody
	}
	return []byte{}
}

func login() {
	fmt.Printf("username: %s, password: %s\n",username,password)
	if username != "" {
		url = fmt.Sprintf("http://localhost%s/TodosV2", lib.PortNo)
		body,_ := json.Marshal(router.UserPutBody{
			Username: username,
			Password: password,
		})

		res:= sendAndReceive(http.MethodPost,body, "http://localhost:8080/Users") //TODO
		if len(res) >0 {
			jwt= string(res)
		} else {
			fmt.Println("Login failed")
			os.Exit(1) 
		}
	} else {
		url = fmt.Sprintf("http://localhost%s/Todos", lib.PortNo)
	}
}
