package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"messsage"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(messsage string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: messsage,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
