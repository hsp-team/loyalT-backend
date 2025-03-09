package errorz

import "errors"

var (
	ErrCoinProgramAlreadyExists = errors.New("coin program already exists")
	ErrCoinProgramNotFound      = errors.New("coin program not found")
)
