package expr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrMockError = errors.New("mock error")

func Test_Result_Type(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Result{Error: ErrNotString}.Type(), TypeError)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Result{Value: "string"}.Type(), TypeString)
	})

	t.Run("Int", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Result{Value: 1}.Type(), TypeInt)
	})

	t.Run("Float", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Result{Value: 1.1}.Type(), TypeFloat)
	})

	t.Run("Bool", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Result{Value: true}.Type(), TypeBool)
	})

	t.Run("Unknown", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Result{}.Type(), TypeUnknown)
	})
}

func Test_Result_As_String(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		expect := "string"
		actual, err := Result{Value: expect}.String()
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("Not_Of_Type", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Value: 1}.String()
		assert.ErrorIs(t, ErrNotString, err)
	})

	t.Run("Eval_Error", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Error: ErrMockError}.String()
		assert.ErrorIs(t, ErrMockError, err)
	})

	t.Run("Must_Ok", func(t *testing.T) {
		t.Parallel()
		expect := "string"
		assert.Equal(t, expect, Result{Value: expect}.MustString())
	})

	t.Run("Must_Error", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithError(t, ErrNotString.Error(), func() {
			Result{Value: 1}.MustString()
		})
	})
}

func Test_Result_As_Int(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		expect := 1
		actual, err := Result{Value: expect}.Int()
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("Not_Of_Type", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Value: "string"}.Int()
		assert.ErrorIs(t, ErrNotInt, err)
	})

	t.Run("Eval_Error", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Error: ErrMockError}.Int()
		assert.ErrorIs(t, ErrMockError, err)
	})

	t.Run("Must_Ok", func(t *testing.T) {
		t.Parallel()
		expect := 1
		assert.Equal(t, expect, Result{Value: expect}.MustInt())
	})

	t.Run("Must_Error", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithError(t, ErrNotInt.Error(), func() {
			Result{Value: "string"}.MustInt()
		})
	})
}

func Test_Result_As_Float(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		expect := 1.1
		actual, err := Result{Value: expect}.Float()
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("Not_Of_Type", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Value: "string"}.Float()
		assert.ErrorIs(t, ErrNotFloat, err)
	})

	t.Run("Eval_Error", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Error: ErrMockError}.Float()
		assert.ErrorIs(t, ErrMockError, err)
	})

	t.Run("Must_Ok", func(t *testing.T) {
		t.Parallel()
		expect := 1.1
		assert.Equal(t, expect, Result{Value: expect}.MustFloat())
	})

	t.Run("Must_Error", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithError(t, ErrNotFloat.Error(), func() {
			Result{Value: "string"}.MustFloat()
		})
	})
}

func Test_Result_As_Bool(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		expect := true
		actual, err := Result{Value: expect}.Bool()
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("Not_Of_Type", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Value: "string"}.Bool()
		assert.ErrorIs(t, ErrNotBool, err)
	})

	t.Run("Eval_Error", func(t *testing.T) {
		t.Parallel()
		_, err := Result{Error: ErrMockError}.Bool()
		assert.ErrorIs(t, ErrMockError, err)
	})

	t.Run("Must_Ok", func(t *testing.T) {
		t.Parallel()
		expect := true
		assert.Equal(t, expect, Result{Value: expect}.MustBool())
	})

	t.Run("Must_Error", func(t *testing.T) {
		t.Parallel()
		assert.PanicsWithError(t, ErrNotBool.Error(), func() {
			Result{Value: "string"}.MustBool()
		})
	})
}
