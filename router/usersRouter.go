package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"github.com/PhilAldridge/TODO-GO/users"
	"github.com/golang-jwt/jwt/v5"
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

func (h *UserApiHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	var body UserPutBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if body.Username == "" || body.Password == "" || err != nil {
		http.Error(w, "You must provide a username and password", http.StatusBadRequest)
		return
	}
	uuid, err := h.users.CreateUser(body.Username, body.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(uuid.String()))
}

func (h *UserApiHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(token))

}


