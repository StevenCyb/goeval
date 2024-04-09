package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Eval(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name        string
		expression  string
		expressionA []interface{}
		result      Result
		expectPanic bool
	}{
		{name: "Number", expression: "  2 ", result: Result{Value: float64(2)}},
		{name: "Boolean", expression: " true ", result: Result{Value: true}},
		{name: "String_Double_Quote", expression: ` "hello" `, result: Result{Value: "hello"}},
		{name: "String_Single_Quote", expression: ` 'world' `, result: Result{Value: "world"}},
		{name: "", expression: ` 'hello' + "world" `, result: Result{Value: "hello" + "world"}},
		{name: "", expression: ` 'hello' + "world" +'!' `, result: Result{Value: "hello" + "world" + "!"}},
		{name: "", expression: "1+  2 ", result: Result{Value: float64(3)}},
		{name: "", expression: "%d+%d ", expressionA: []interface{}{2, 2}, result: Result{Value: float64(4)}},
		{name: "", expression: "6/  2 ", result: Result{Value: float64(3.0)}},
		{name: "", expression: "2+2*3 ", result: Result{Value: float64(8)}},
		{name: "", expression: "1+2*3 -1", result: Result{Value: float64(6)}},
		{name: "", expression: "1+2*3 -1*2", result: Result{Value: float64(5)}},
		{name: "", expression: " true  && true ", result: Result{Value: true}},
		{name: "", expression: " true  && false ", result: Result{Value: false}},
		{name: "", expression: " false  || true ", result: Result{Value: true}},
		{name: "", expression: " false  || false &&true ", result: Result{Value: false}},
		{name: "", expression: " false  || false && true || false ", result: Result{Value: false}},
	}
	/*
	   1>3
	   1<2<3
	   true==true
	   a==a!=b
	   true!=false
	   "a"+"b"=="ab"



	   (1+2)*3
	   ("a"+"b")+"c"
	   ("a"!="b")==true
	*/

	for _, tc := range tcs {
		tcRef := tc

		t.Run(tcRef.name, func(t *testing.T) {
			t.Parallel()

			if tcRef.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic")
					}
				}()
			}

			result := Eval(tcRef.expression, tcRef.expressionA...)
			assert.Equal(t, tcRef.result, result)
		})
	}
}
