package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Eval(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		expression  string
		expressionA []interface{}
		result      Result
		expectPanic bool
	}{
		{expression: "1", result: Result{Value: float64(1)}},
		{expression: "  2 ", result: Result{Value: float64(2)}},
		{expression: " true ", result: Result{Value: true}},
		{expression: ` "hello" `, result: Result{Value: "hello"}},
		{expression: ` 'world' `, result: Result{Value: "world"}},
		{expression: ` 'hello' + "world" `, result: Result{Value: "hello" + "world"}},
		{expression: ` 'hello' + "world" +'!' `, result: Result{Value: "hello" + "world" + "!"}},
		{expression: "1+  2 ", result: Result{Value: float64(3)}},
		{expression: "%d+%d ", expressionA: []interface{}{2, 2}, result: Result{Value: float64(4)}},
		{expression: "6/  2 ", result: Result{Value: float64(3.0)}},
		{expression: "2+2*3 ", result: Result{Value: float64(8)}},
		{expression: "1+2*3 -1", result: Result{Value: float64(6)}},
		{expression: "1+2*3 -1*2", result: Result{Value: float64(5)}},
		{expression: " true  && true ", result: Result{Value: true}},
		{expression: " true  && false ", result: Result{Value: false}},
		{expression: " false  || true ", result: Result{Value: true}},
		{expression: " false  || false &&true ", result: Result{Value: false}},
		{expression: " false  || false && true || false ", result: Result{Value: false}},
	}

	for _, tc := range tcs {
		tcRef := tc

		t.Run(tcRef.expression, func(t *testing.T) {
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
