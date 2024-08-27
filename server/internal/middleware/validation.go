// File: internal/middleware/validation.go
package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/JettZgg/LineUp/internal/utils"
)

func ValidateJSONBody(next http.HandlerFunc, v interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			utils.NewAppError(err, "Invalid request body", http.StatusBadRequest).LogAndRespond(w)
			return
		}
		next.ServeHTTP(w, r)
	}
}