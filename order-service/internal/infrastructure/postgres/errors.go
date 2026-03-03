package postgres

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order with this ID already exists")
	ErrNotFound           = errors.New("order with that id doenot exist")
)
