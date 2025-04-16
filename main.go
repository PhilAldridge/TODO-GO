package main

import (
	"fmt"
	"os"

	"github.com/PhilAldridge/TODO-GO/cli"
)

func main() {
	cmd := cli.NewCmd(os.Stdout)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
