package domain

import "errors"

var (
	ErrDuplicateEmail = errors.New("email already in use")
	ErrNotFound       = errors.New("not found")
	ErrAuthentication = errors.New("authentication err")
	ErrAuthorization  = errors.New("authorization err")
	ErrDuplicateKey   = errors.New("duplicate key err")
)
