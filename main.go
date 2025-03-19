package main

import (
	"fmt"
	"os"

	"github.com/PhilAldridge/TODO-GO/cli"
	"github.com/PhilAldridge/TODO-GO/store"
)

func main() {
	var myStore store.Store = store.NewTodoStore()

	cmd := cli.NewCmd(os.Stdout, myStore)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
