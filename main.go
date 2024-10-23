package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	todoCli(os.Stdin, os.Stdout)

}

func todoCli(reader io.Reader, writer io.Writer) {
	fmt.Fprintln(writer, "Todo List:")
	input := bufio.NewScanner(reader)
	for {
		input.Scan()
		text := input.Text()
		fmt.Fprintln(writer, text)
		if input.Text() == "exit" {
			fmt.Fprintln(writer, "Exiting App")
			break
		}
	}
}
