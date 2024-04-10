package errs

import (
	"fmt"
)

const errUnexpectedTokenTypeMessage = "Unexpected token: \"%s\", expected: \"%s\""

// UnexpectedTokenTypeError.TokenType is an error
// type for unexpected token type.
type UnexpectedTokenTypeError struct {
	actual   string
	expected string
}

// Error returns the error message text.
func (err UnexpectedTokenTypeError) Error() string {
	return fmt.Sprintf(errUnexpectedTokenTypeMessage,
		err.actual, err.expected)
}

// NewErrUnexpectedTokenType cerate a new error.
func NewErrUnexpectedTokenType(actual, expected string) UnexpectedTokenTypeError {
	return UnexpectedTokenTypeError{
		actual:   actual,
		expected: expected,
	}
}
