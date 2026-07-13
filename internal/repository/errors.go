package repository

import "errors"

// Sentinel errors returned by repositories.
var (
	ErrNotFound          = errors.New("not found")
	ErrDuplicate         = errors.New("duplicate entry")
	ErrInsufficientStock = errors.New("insufficient stock")
)
