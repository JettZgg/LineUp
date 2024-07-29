// File: internal/utils/errors.go
package utils

import (
	"log"
	"net/http"
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func NewAppError(err error, message string, code int) AppError {
	return AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func (ae AppError) LogAndRespond(w http.ResponseWriter) {
	log.Printf("Error: %v", ae.Err)
	http.Error(w, ae.Message, ae.Code)
}
