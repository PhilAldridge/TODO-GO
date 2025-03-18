package main

import (
	"fmt"
	"os"

	"github.com/PhilAldridge/TODO-GO/cli"
	"github.com/PhilAldridge/TODO-GO/store"
)

func main() {
	store := store.NewTodoList()

	cmd := cli.NewCmd(os.Stdout,store)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
