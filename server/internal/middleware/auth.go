package middleware

import (
	"net/http"
	"strings"

	"github.com/JettZgg/LineUp/internal/auth"
	"github.com/JettZgg/LineUp/internal/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.NewAppError(nil, "Missing authorization header", http.StatusUnauthorized).LogAndRespond(w)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			utils.NewAppError(nil, "Invalid authorization header", http.StatusUnauthorized).LogAndRespond(w)
			return
		}

		username, err := auth.ValidateToken(bearerToken[1])
		if err != nil {
			utils.NewAppError(err, "Invalid token", http.StatusUnauthorized).LogAndRespond(w)
			return
		}

		r.Header.Set("Username", username)
		next.ServeHTTP(w, r)
	}
}