package router

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func internalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func marshalAndWrite[T any](w http.ResponseWriter, data T) {
	resp, err := json.Marshal(data)
	if err != nil {
		internalServerErrorHandler(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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

func ServeTemplate(path string, data interface{}, w http.ResponseWriter) {
	funcMap := template.FuncMap{
		"formatUKDate": func(t time.Time) string {
			return t.Format("02/01/2006") // DD/MM/YYYY
		},
	}

	tmpl,err :=  template.New("list.html").Funcs(funcMap).ParseFiles(path)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

