package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorAtPositionError(t *testing.T) {
	t.Parallel()

	pos := 1
	key := "b"
	require.Equal(t,
		fmt.Sprintf(errErrorAtPosition, key, pos),
		NewErrErrorAtPosition(fmt.Errorf(key), pos).Error(),
	)
}
