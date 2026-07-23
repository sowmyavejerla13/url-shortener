package apperrors

import "errors"

var (
	ErrURLNotFound        = errors.New("url not found")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidURL         = errors.New("invalid url")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
