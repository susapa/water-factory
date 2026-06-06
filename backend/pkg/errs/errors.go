package errs

import "net/http"

type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *APIError) Error() string { return e.Message }

func New(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

func BadRequest(msg string) *APIError  { return New(http.StatusBadRequest, msg) }
func Unauthorized(msg string) *APIError { return New(http.StatusUnauthorized, msg) }
func Forbidden(msg string) *APIError    { return New(http.StatusForbidden, msg) }
func NotFound(msg string) *APIError     { return New(http.StatusNotFound, msg) }
func Conflict(msg string) *APIError     { return New(http.StatusConflict, msg) }
func Internal(msg string) *APIError     { return New(http.StatusInternalServerError, msg) }
