package apperrors

import "errors"

var (
	ErrNotFound     = errors.New("Not found")
	ErrAlreadyExist = errors.New("Already exist")
	ErrAccessDenied = errors.New("Acess denied")
)
