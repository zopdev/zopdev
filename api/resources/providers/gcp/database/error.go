package sql

import "net/http"

type ErrConflict struct {
	Message string `json:"message"`
}

func (e *ErrConflict) Error() string {
	return e.Message
}

func (e *ErrConflict) StatusCode() int {
	return http.StatusConflict
}

type InternalServerError struct{}

func (e *InternalServerError) Error() string {
	return "Internal server error!"
}
