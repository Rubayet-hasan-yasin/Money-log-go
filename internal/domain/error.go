package domain

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyInUse    = errors.New("email already in use")
	ErrUsernameAlreadyInUse = errors.New("username already in use")
	ErrInvalidToken         = errors.New("invalid or expired token")
)


