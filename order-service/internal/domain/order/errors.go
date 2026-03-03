package order

import "errors"

var (
	ErrInvalidMail    = errors.New("that mail is invalid")
	ErrNoItems        = errors.New("order must contain items or at least one of them")
	ErrEmptyID        = errors.New("order id cannot be empty")
	ErrInvalidPayment = errors.New("payment amount cannot be zero or below it")
)
