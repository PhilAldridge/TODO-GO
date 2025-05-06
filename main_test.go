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

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/router"
	"github.com/PhilAldridge/TODO-GO/store"
)

func TestMain(m *testing.M) {
	lib.LoadConfig(".env.test")
	os.Exit(m.Run())
}

func TestConcurrentPutAndGet(t *testing.T) {
	todoStore, userStore := store.NewSQLStore()
	server := SetupServer(todoStore,userStore)

	for i := 0; i < 50; i++ {
		t.Run(fmt.Sprintf("ParallelTestV1-%d", i), func(t *testing.T) {
			t.Parallel()

			payload := router.V1PutBody{
				Label:    "test" + strconv.Itoa(i),
				Deadline: "2025-01-01",
			}

			resp := sendRequest(server, http.MethodPut, "/Todos", payload, map[string]string{
				"Content-Type": "application/json",
			})
			mustStatusOK(t, resp, "PUT")

			resp = sendRequest(server, http.MethodGet, "/Todos", nil, nil)
			mustStatusOK(t, resp, "GET")
		})
	}
}

func TestConcurrentMultipleUsers(t *testing.T) {
	todoStore, userStore := store.NewSQLStore()
	server := SetupServer(todoStore,userStore)

	for i := 0; i < 50; i++ {
		t.Run(fmt.Sprintf("ParallelTestV2-%d", i), func(t *testing.T) {
			t.Parallel()

			user := router.UserPutBody{
				Username: "user" + strconv.Itoa(i),
				Password: "password" + strconv.Itoa(i),
			}

			sendRequest(server, http.MethodPut, "/Users", user, map[string]string{
				"Content-Type": "application/json",
			})

			resp := sendRequest(server, http.MethodPost, "/Users", user, map[string]string{
				"Content-Type": "application/json",
			})
			mustStatusOK(t, resp, "Login")

			jwt := resp.Body.String()

			todo := router.V1PutBody{
				Label:    "test" + strconv.Itoa(i),
				Deadline: "2025-01-01",
			}

			resp = sendRequest(server, http.MethodPut, "/TodosV2", todo, map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + jwt,
			})
			mustStatusOK(t, resp, "PUT with JWT")
		})
	}
}

func sendRequest(server http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var b io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		b = bytes.NewBuffer(jsonBytes)
	}

	req := httptest.NewRequest(method, path, b)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	return w
}

func mustStatusOK(t *testing.T, resp *httptest.ResponseRecorder, label string) {
	t.Helper()
	if resp.Code != http.StatusOK {
		t.Errorf("%s: unexpected status: %d", label, resp.Code)
	}
}
