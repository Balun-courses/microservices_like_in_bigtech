package models

import "errors"

var (
	// ErrAlreadyExists - error already exists
	ErrAlreadyExists = errors.New("already exists")
	// ErrUnimplemented - error unimplemented
	ErrUnimplemented = errors.New("unimplemented")
)
