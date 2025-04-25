package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/PhilAldridge/TODO-GO/auth"
	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/logging"
	"github.com/PhilAldridge/TODO-GO/router"
	"github.com/PhilAldridge/TODO-GO/store"
	"github.com/PhilAldridge/TODO-GO/users"
)

var (
	mode = flag.String("mode", "mem", "set the mode the application should run in (mem, json, sql)")
)

func main() {

	lib.LoadConfig(".env")

	flag.Parse()

	var todoStore store.Store
	var usersStore users.Users
	switch *mode {
	case "mem":
		todoStore = store.NewInMemoryTodoStore()
		usersStore = users.NewInMemoryUsersStore()
	case "json":
		todoStore = &store.JSONStore{}
		usersStore = &users.JSONUsers{}
	default:
		log.Fatal("valid modes: json, mem")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	fmt.Println("Server listening on :8080")
	mux := http.NewServeMux()
	v1api := router.NewV1ApiHandler(todoStore)
	usersapi := router.NewUserApiHandler(usersStore)
	v2api := router.NewV2ApiHandler(todoStore)

	mux.Handle("/Todos/", &v1api)
	mux.Handle("/Todos", &v1api)
	mux.Handle("/Users/", &usersapi)
	mux.Handle("/Users", &usersapi)
	mux.Handle("/TodosV2/", auth.JWTMiddleware(&v2api))
	mux.Handle("/TodosV2", auth.JWTMiddleware(&v2api))
	wrapped := logging.WithTraceIDAndLogger(
		logging.LoggingMiddleware(mux),
	)

	http.ListenAndServe(lib.PortNo, wrapped)
}
