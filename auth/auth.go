package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//get authorisation token from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//check token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return lib.JwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//add username from token to context
		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		ctx := context.WithValue(r.Context(), contextKey("username"), username)

		// Token is valid — continue
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
