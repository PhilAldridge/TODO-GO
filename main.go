package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/PhilAldridge/TODO-GO/logging"
	"github.com/PhilAldridge/TODO-GO/router"
	"github.com/PhilAldridge/TODO-GO/store"
)

var (
	mode = flag.String("mode","mem","set the mode the application should run in (mem, json, sql)")
)

func main() {
	// cmd := cli.NewCmd(os.Stdout)
	// if err := cmd.Execute(); err != nil {
	// 	fmt.Fprintln(os.Stdout, err)
	// 	os.Exit(1)
	// }
	flag.Parse()

	var todoStore store.Store
	switch *mode {
	case "mem":
		todoStore = store.NewInMemoryTodoStore()
	case "json":
		todoStore = &store.JSONStore{}
	default:
		log.Fatal("valid modes: json, mem")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    slog.SetDefault(logger)

	fmt.Println("Server listening on :8080")
	mux := http.NewServeMux()
	api:= router.NewApiHandler(todoStore)

	mux.Handle("/Todos/", &api)
	mux.Handle("/Todos", &api)
	wrapped := logging.WithTraceIDAndLogger(
		logging.LoggingMiddleware(mux),
	)

    http.ListenAndServe(":8080", wrapped)
}
