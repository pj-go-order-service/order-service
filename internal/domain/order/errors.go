package order

import "errors"

var (
	ErrEmptyOrder   = errors.New("order must contain at least one item")
	ErrInvalidState = errors.New("invalid order state")
)
