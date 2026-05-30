package middleware

import (
	"net/http"
	"strings"
)

func ValidateToken(ExpectedToken string) func (http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "Missing Authorization", http.StatusUnauthorized)
				return 
			}

			token, ok := strings.CutPrefix(authHeader, "Bearer ")

			if !ok || token == "" {
				http.Error(w, "Missing Token Value", http.StatusUnauthorized)
				return 
			}

			if token != ExpectedToken {
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
			}

			next.ServeHTTP(w, r)
		})
	}
}