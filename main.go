package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/PhilAldridge/TODO-GO/store"
)

// type command struct {
// 	commandText string
// 	commandLogic func()
// }

// commandList:= []command{
// 	{}
// }

func main() {
	todoCli(os.Stdin, os.Stdout)

}

func todoCli(reader io.Reader, writer io.Writer) {
	fmt.Fprintln(writer, "***Welcome to my todo app***\nType a command or type 'help' to list commands:")
	input := bufio.NewScanner(reader)
	todoList := store.NewTodoList()
	for {
		input.Scan()
		switch input.Text() {
		case "exit":
			fmt.Fprintln(writer, "Exiting App")
			return
		case "add":
			fmt.Fprintln(writer, "Write your todo")
			input.Scan()
			label := input.Text()
			fmt.Fprintln(writer, "Write your deadline (in dd/mm/yyyy format)")
			input.Scan()
			deadline, _ := time.Parse("01/02/2006", input.Text())
			todoList.AddTodo(label, "", deadline)
		case "list":
			todos := todoList.ListTodos()
			for _, todo := range todos {
				fmt.Fprintln(writer, todo.Label)
			}
		case "help":
			fmt.Fprintln(writer, "add\nlist\nexit\nType a command:")
		}
	}
}
