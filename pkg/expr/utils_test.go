package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertFloat(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name   string
		value  interface{}
		expect float64
	}{
		{name: "Float64", value: float64(1), expect: 1},
		{name: "String", value: "hello", expect: 5},
		{name: "Bool_True", value: true, expect: 1},
		{name: "Bool_False", value: false, expect: 0},
	}

	for _, tc := range tcs {
		tcRef := tc

		t.Run(tcRef.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tcRef.expect, convertFloat(tcRef.value))
		})
	}
}
