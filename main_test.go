package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

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
	store:= &store.JSONStore{}
	api:= router.NewApiHandler(store)
	server.Handle("/Todos/", &api)
	server.Handle("/Todos", &api)

	for i := 0; i < 10; i++ {
		t.Run("ParallelTest", func(t *testing.T) {
			t.Parallel()

			payload := router.PutBody{
				Label: "test"+strconv.Itoa(i),
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