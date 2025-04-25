package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserApiHandler struct {
	users users.Users
}

type Claims struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	jwt.RegisteredClaims
}

func NewUserApiHandler(users users.Users) UserApiHandler {
	return UserApiHandler{users: users}
}

func (h *UserApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPut:
		h.handlePut(w, r)
	case r.Method == http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
}

func (h *UserApiHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	var body UserPutBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if body.Username == "" || body.Password == "" || err != nil {
		http.Error(w, "You must provide a username and password", http.StatusBadRequest)
		return
	}
	uuid, err := h.users.CreateUser(body.Username, body.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusConflict) //TODO
		return
	}
	w.Write([]byte(uuid.String()))
}

func (h *UserApiHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var body UserPutBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if body.Username == "" || body.Password == "" || err != nil {
		http.Error(w, "You must provide a username and password", http.StatusBadRequest)
		return
	}
	uuid, err := h.users.Login(body.Username, body.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusUnauthorized)
		return
	}
	token, err := createToken(body.Username, uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusFailedDependency)
		return
	}
	fmt.Println(token)
	w.Write([]byte(token))

}

func createToken(username string, id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"id":       id.String(),
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(lib.JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
