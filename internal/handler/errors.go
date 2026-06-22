package handler

import (
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

var (
	ErrInternalServerError = HTTPError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrBadRequest          = HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}
	ErrUnauthorized        = HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden           = HTTPError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrNotFound            = HTTPError{Code: http.StatusNotFound, Message: "Not found"}
)
