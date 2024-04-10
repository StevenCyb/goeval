package expr

import (
	"fmt"
	"testing"

	"github.com/StevenCyb/goeval/pkg/errs"
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
		{name: "Empty", expression: "", result: Result{Value: nil, Error: errs.ErrEmptyExpression}},
		{name: "Number", expression: "  2 ", result: Result{Value: float64(2)}},
		{name: "Boolean", expression: " true ", result: Result{Value: true}},
		{name: "String_Double_Quote", expression: ` "hello" `, result: Result{Value: "hello"}},
		{name: "String_Single_Quote", expression: ` 'world' `, result: Result{Value: "world"}},
		{name: "Arithmetic_String", expression: ` 'hello' - "world" `, result: Result{Value: float64(0)}},
		{name: "Chained_Arithmetic_String", expression: ` 'hello' + "world" +'!' `, result: Result{Value: float64(11)}},
		{name: "Add", expression: "1+  2 ", result: Result{Value: float64(3)}},
		{name: "Add_Formatted", expression: "%d+%d ", expressionA: []interface{}{2, 2}, result: Result{Value: float64(4)}},
		{name: "Divide", expression: "6/  2 ", result: Result{Value: float64(3.0)}},
		{name: "Chained_Calc", expression: "2+2*3 ", result: Result{Value: float64(8)}},
		{name: "Chained_Calc_Precedence", expression: "1+2*3 -1", result: Result{Value: float64(6)}},
		{name: "Chained_Calc_Precedence2", expression: "1+2*3 -1*2", result: Result{Value: float64(5)}},
		{name: "Logical_And", expression: " true  && true ", result: Result{Value: true}},
		{name: "Logical_Or", expression: " false  || true ", result: Result{Value: true}},
		{name: "Chained_Logical", expression: " false  || false &&true ", result: Result{Value: false}},
		{name: "Longer_Chained_Logical_False", expression: " false  || false && true || false ", result: Result{Value: false}},
		{name: "Comparison_Number", expression: " 1 > 3 ", result: Result{Value: false}},
		{name: "Chained_Comparison_Number", expression: " 1+5 > 3 ", result: Result{Value: true}},
		{name: "Comparison_String_Number", expression: " 3 <= 'hello' ", result: Result{Value: true}},
		{name: "Comparison_Bool", expression: " true == true ", result: Result{Value: true}},
		{name: "Added_String_Comparison", expression: ` "a"+"b"==2`, result: Result{Value: true}},
		{name: "Chained_Equal_Bool", expression: " 'a'=='a'!='b' ", result: Result{Value: false}},
		{name: "Simple_Context", expression: " (1+2) ", result: Result{Value: 3}},
		// {name: "Chained_Simple_Context", expression: " (1+2)*2 ", result: Result{Value: 6}},
	}
	/*
		TODO context
		   (1+2)*3
		   ("a"+"b")+"c"
		   ("a"!="b")==true
	*/

	for _, tc := range tcs {
		tcRef := tc

		t.Run(tcRef.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if r := recover(); r == nil {
					if tcRef.expectPanic {
						t.Errorf("expected panic")
					}
				} else {
					t.Errorf("unexpected panic on '%s': %v", tcRef.name, r)
				}
			}()

			result := Eval(tcRef.expression, tcRef.expressionA...)
			assert.Equal(t, tcRef.result, result, fmt.Sprintf(tc.expression, tc.expressionA...))
		})
	}
}
