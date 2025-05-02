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
	_ "github.com/lib/pq"
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
	case "sql":
		todoStore,usersStore = store.NewSQLStore()
		defer todoStore.Close()
	default:
		log.Fatal("valid modes: json, mem,sql")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	fmt.Println("Server listening on :8080")
	mux := http.NewServeMux()
	v1api := router.NewV1ApiHandler(todoStore)
	usersapi := router.NewUserApiHandler(usersStore)
	v2api := router.NewV2ApiHandler(todoStore)
	
	mux.HandleFunc("GET /Todos",v1api.HandleGet)
	mux.HandleFunc("PATCH /Todos",v1api.HandlePatch)
	mux.HandleFunc("PUT /Todos",v1api.HandlePut)
	mux.HandleFunc("DELETE /Todos",v1api.HandleDelete)
	mux.HandleFunc("PUT /Users",usersapi.HandlePut)
	mux.HandleFunc("POST /Users",usersapi.HandlePost)
	mux.HandleFunc("GET /TodosV2",auth.JWTMiddleware(v2api.HandleGet))
	mux.HandleFunc("PATCH /TodosV2",auth.JWTMiddleware(v2api.HandlePatch))
	mux.HandleFunc("PUT /TodosV2",auth.JWTMiddleware(v2api.HandlePut))
	mux.HandleFunc("DELETE /TodosV2",auth.JWTMiddleware(v2api.HandleDelete))

	wrapped := logging.WithTraceIDAndLogger(
		logging.LoggingMiddleware(mux),
	)

	http.ListenAndServe(lib.PortNo, wrapped)
}
