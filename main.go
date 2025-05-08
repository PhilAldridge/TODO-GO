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
	"syscall"
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

	todoStore, usersStore:= SetupStores(*mode)
	srv,actor := SetupServer(todoStore,usersStore)

	
	// Channel to listen for shutdown signals
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server listening on :8080")
		http.ListenAndServe(lib.PortNo, srv.Handler)
	}()
	

	<-stopCh
	shutdown(srv,actor)
}

func SetupStores(mode string) (store.Store, users.Users) {
	var todoStore store.Store
	var usersStore users.Users
	switch mode {
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
	return todoStore,usersStore
}

func SetupServer(todoStore store.Store, userStore users.Users) (*http.Server,*store.StoreActor) {
	mux := http.NewServeMux()
	usersapi := router.NewUserApiHandler(userStore)
	actor:=store.StartStoreActor(todoStore)
	v2api := router.NewV2ApiHandler(actor)
	fs:= http.FileServer(http.Dir("./static"))
	
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	wrapped := logging.WithTraceIDAndLogger(
		logging.LoggingMiddleware(mux),
	)
	srv:= &http.Server{
		Addr: lib.PortNo,
		Handler: wrapped,
	}

	return srv, actor
}

func shutdown(srv *http.Server, actor *store.StoreActor) {
	fmt.Println("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		os.Exit(1)
	}

	actor.Stop()
	log.Println("Server shutdown complete")
}

