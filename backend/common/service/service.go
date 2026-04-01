package service

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrBadRequest    = errors.New("bad request")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal")
	ErrNotFound      = errors.New("record not found")
)
