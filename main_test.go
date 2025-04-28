package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/PhilAldridge/TODO-GO/auth"
	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/router"
	"github.com/PhilAldridge/TODO-GO/store"
)

func TestMain(m *testing.M) {
	lib.LoadConfig(".env.test")
	os.Exit(m.Run())
}

func TestConcurrentPutAndGet(t *testing.T) {
	server := http.NewServeMux()
	store := &store.JSONStore{}
	api := router.NewV1ApiHandler(store)
	server.Handle("/Todos/", &api)
	server.Handle("/Todos", &api)

	for i := 0; i < 50; i++ {
		t.Run("ParallelTestV1", func(t *testing.T) {
			t.Parallel()

			payload := router.V1PutBody{
				Label:    "test" + strconv.Itoa(i),
				Deadline: "2025-01-01",
			}
			b, _ := json.Marshal(payload)

			req := httptest.NewRequest(http.MethodPut, "/Todos", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("PUT: unexpected status: %d", w.Code)
			}

			req = httptest.NewRequest(http.MethodGet, "/Todos", nil)
			w = httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("GET: unexpected status: %d", w.Code)
			}
		})
	}
}

func TestConcurrentPutAndGetMultipleUsers(t *testing.T) {
	server := http.NewServeMux()
	todoStore, userStore := store.NewSQLStore()
	todoApi := router.NewV2ApiHandler(todoStore)
	server.Handle("/TodosV2/", auth.JWTMiddleware(&todoApi))
	server.Handle("/TodosV2", auth.JWTMiddleware(&todoApi))

	userApi := router.NewUserApiHandler(userStore)
	server.Handle("/Users/", &userApi)
	server.Handle("/Users", &userApi)

	for i := 2; i < 50; i++ {
		t.Run("ParallelTestV2", func(t *testing.T) {
			t.Parallel()

			payload := router.UserPutBody{
				Username: "user" + strconv.Itoa(i),
				Password: "password" + strconv.Itoa(i),
			}
			b, _ := json.Marshal(payload)

			//Create user
			req:= httptest.NewRequest(http.MethodPut, "/Users", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			w:= httptest.NewRecorder()
			server.ServeHTTP(w, req)

			//Login
			req= httptest.NewRequest(http.MethodPost, "/Users", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			w= httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Login: unexpected status: %d", w.Code)
			}

			resBody, err := io.ReadAll(w.Body)
			if err != nil {
				t.Errorf("client: could not read response body: %s\n", err)

			}

			//Use JWT to authorise PUT method in V2API
			jwt:= string(resBody)
			fmt.Println(jwt)
			V2payload:= router.V1PutBody{
				Label:    "test" + strconv.Itoa(i),
				Deadline: "2025-01-01",
			}
			b,_= json.Marshal(V2payload)

			req= httptest.NewRequest(http.MethodPut, "/TodosV2", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization","Bearer "+jwt)
			w = httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Put: unexpected status: %d", w.Code)
			}

		})
	}
}
