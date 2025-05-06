package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	mux := SetupServer(todoStore,usersStore)

	wrapped := logging.WithTraceIDAndLogger(
		logging.LoggingMiddleware(mux),
	)
	srv:= &http.Server{
		Addr: lib.PortNo,
		Handler: mux,
	}
	idleConnsClosed:= shutdownChannel(srv)
	http.ListenAndServe(lib.PortNo, wrapped)

	<-idleConnsClosed
	log.Println("Server shutdown complete")
}

func SetupServer(todoStore store.Store, userStore users.Users) http.Handler {
	mux := http.NewServeMux()
	//v1api := router.NewV1ApiHandler(todoStore)
	usersapi := router.NewUserApiHandler(userStore)
	v2api := router.NewV2ApiHandler(todoStore)
	fs:= http.FileServer(http.Dir("./static"))

	mux.Handle("/",fs)
	mux.HandleFunc("GET /Todos",v2api.HandleGet)
	mux.HandleFunc("PATCH /Todos",v2api.HandlePatch)
	mux.HandleFunc("PUT /Todos",v2api.HandlePut)
	mux.HandleFunc("DELETE /Todos",v2api.HandleDelete)
	mux.HandleFunc("PUT /Users",usersapi.HandlePut)
	mux.HandleFunc("POST /Users",usersapi.HandlePost)
	mux.HandleFunc("GET /TodosV2",auth.JWTMiddleware(v2api.HandleGet))
	mux.HandleFunc("PATCH /TodosV2",auth.JWTMiddleware(v2api.HandlePatch))
	mux.HandleFunc("PUT /TodosV2",auth.JWTMiddleware(v2api.HandlePut))
	mux.HandleFunc("DELETE /TodosV2",auth.JWTMiddleware(v2api.HandleDelete))
	mux.HandleFunc("GET /List",v2api.HandleList)

	return mux
}

func shutdownChannel(srv *http.Server) chan struct{} {
	// Channel to signal when shutdown is complete
	idleConnsClosed := make(chan struct{})

	// Handle interrupt signal for graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutdown signal received")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
			os.Exit(1)
		}
		close(idleConnsClosed)
	}()

	return idleConnsClosed
}

