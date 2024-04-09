package expr

import (
	"fmt"
	"strings"

	"github.com/StevenCyb/goeval/pkg/errs"
)

func Eval(format string, a ...any) Result {
	expression := strings.TrimSpace(fmt.Sprintf(format, a...))
	if expression == "" {
		return Result{
			Error: errs.ErrEmptyExpression,
		}
	}

	return newParser(expression).Parse()
}
