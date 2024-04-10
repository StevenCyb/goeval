package errs

import "errors"

var (
	ErrEmptyExpression = errors.New("empty expression")
	ErrDivisionByZero  = errors.New("division by zero")
)
