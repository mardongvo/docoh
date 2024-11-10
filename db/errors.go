package db

import "errors"

var (
	ErrDuplicateRule = errors.New("duplicate rule")
	ErrFoundNoRule   = errors.New("found no rule")
)
