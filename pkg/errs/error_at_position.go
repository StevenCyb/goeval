package errs

import "fmt"

const errErrorAtPosition = "%s, at position %d"

// ErrorAtPositionError is an error.
type ErrorAtPositionError struct {
	err      error
	position int
}

// Error returns the error message text.
func (err ErrorAtPositionError) Error() string {
	return fmt.Sprintf(errErrorAtPosition, err.err.Error(), err.position)
}

// NewErrErrorAtPosition cerate a new error.
func NewErrErrorAtPosition(err error, position int) ErrorAtPositionError {
	return ErrorAtPositionError{
		err:      err,
		position: position,
	}
}
