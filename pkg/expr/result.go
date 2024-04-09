package expr

import "errors"

var (
	ErrNotString = errors.New("value is not a string")
	ErrNotInt    = errors.New("value is not an int")
	ErrNotFloat  = errors.New("value is not a float64")
	ErrNotBool   = errors.New("value is not a bool")
)

type Type string

const (
	TypeError   Type = "error"
	TypeString  Type = "string"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeBool    Type = "bool"
	TypeUnknown Type = "unknown"
)

// Result represents the result of an expression evaluation.
type Result struct {
	Error error
	Value interface{}
}

// Type returns the type of the result.
func (r Result) Type() Type {
	if r.Error != nil {
		return TypeError
	}

	switch r.Value.(type) {
	case string:
		return TypeString
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return TypeInt
	case float32, float64:
		return TypeFloat
	case bool:
		return TypeBool
	default:
		return TypeUnknown
	}
}

// String returns the result as string or error if not a string.
func (r Result) String() (string, error) {
	if r.Error != nil {
		return "", r.Error
	}

	value, ok := r.Value.(string)
	if !ok {
		return "", ErrNotString
	}

	return value, nil
}

// MustString returns the result as string or panics if not a string or eval failed.
func (r Result) MustString() string {
	value, err := r.String()
	if err != nil {
		panic(err)
	}

	return value
}

// Int returns the result as int or error if not an int.
func (r Result) Int() (int, error) {
	if r.Error != nil {
		return 0, r.Error
	}

	value, ok := r.Value.(int)
	if !ok {
		return 0, ErrNotInt
	}

	return value, nil
}

// MustInt returns the result as int or panics if not an int or eval failed.
func (r Result) MustInt() int {
	value, err := r.Int()
	if err != nil {
		panic(err)
	}

	return value
}

// Float returns the result as float64 or error if not a float64.
func (r Result) Float() (float64, error) {
	if r.Error != nil {
		return 0, r.Error
	}

	value, ok := r.Value.(float64)
	if !ok {
		return 0, ErrNotFloat
	}

	return value, nil
}

// MustFloat returns the result as float64 or panics if not a float64 or eval failed.
func (r Result) MustFloat() float64 {
	value, err := r.Float()
	if err != nil {
		panic(err)
	}

	return value
}

// Bool returns the result as bool or error if not a bool.
func (r Result) Bool() (bool, error) {
	if r.Error != nil {
		return false, r.Error
	}

	value, ok := r.Value.(bool)
	if !ok {
		return false, ErrNotBool
	}

	return value, nil
}

// MustBool returns the result as bool or panics if not a bool or eval failed.
func (r Result) MustBool() bool {
	value, err := r.Bool()
	if err != nil {
		panic(err)
	}

	return value
}