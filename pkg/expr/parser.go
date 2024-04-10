package expr

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/StevenCyb/goeval/pkg/errs"
	"github.com/StevenCyb/gotokenizer/pkg/tokenizer"
)

var (
	skipType                tokenizer.Type = "SKIP"
	contextStartType        tokenizer.Type = "CONTEXT_START"
	contextEndType          tokenizer.Type = "CONTEXT_END"
	arithmeticOperationType tokenizer.Type = "ARITHMETIC_OPERATION"
	comparisonOperationType tokenizer.Type = "COMPARISON_OPERATION"
	logicalOperationType    tokenizer.Type = "LOGICAL_OPERATION"
	numberType              tokenizer.Type = "NUMBER"
	boolType                tokenizer.Type = "BOOL"
	textType                tokenizer.Type = "TEXT"

	intBase     = 10
	int64Size   = 64
	float64Size = 64
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

// LL(2) parser for the following grammar:
/*
<EXPRESSION>            ::= <NUMBER>
													| <TEXT>
													| <BOOL>
													| <CONTEXT_EXPRESSION>
													| <ARITHMETIC_EXPRESSION>
													| <LOGICAL_EXPRESSION>
													| <COMPARISON_EXPRESSION>

<ARITHMETIC_EXPRESSION> ::= <EXPRESSION> <ARITHMETIC_OPERATION> <EXPRESSION>
<LOGICAL_EXPRESSION> 		::= <EXPRESSION> <LOGICAL_OPERATION> <EXPRESSION>
<COMPARISON_EXPRESSION>	::= <EXPRESSION> <COMPARISON_OPERATION> <EXPRESSION>
<CONTEXT_EXPRESSION>		::= <CONTEXT_START> <EXPRESSION> <CONTEXT_END>
													| <CONTEXT_START> <EXPRESSION> <CONTEXT_END> <ARITHMETIC_OPERATION> <EXPRESSION>
													| <CONTEXT_START> <EXPRESSION> <CONTEXT_END> <LOGICAL_EXPRESSION> <EXPRESSION>
													| <CONTEXT_START> <EXPRESSION> <CONTEXT_END> <COMPARISON_EXPRESSION> <EXPRESSION>

<SKIP>                  ::= ^\s+
<CONTEXT_START>         ::= ^\(
<CONTEXT_END>           ::= ^\)
<LOGICAL_OPERATION>     ::= ^(&&|\|\|)
<COMPARISON_OPERATION>  ::= ^(==|!=|<=?|>=?)
<ARITHMETIC_OPERATION>  ::= ^(\+|-|\*|\/|%)
<NUMBER>                ::= ^\d+(\.\d+)?
<BOOL>                  ::= ^(true|false)
<TEXT>                  ::= ^("[^"]*"|'[^"]*')
*/
type parser struct {
	tokenizer  *tokenizer.Tokenizer
	lookahead1 *tokenizer.Token
	lookahead2 *tokenizer.Token
}

// Create a new LL(2) parser for the given expression.
func newParser(expression string) *parser {
	return &parser{
		tokenizer: tokenizer.New(
			expression,
			skipType,
			[]*tokenizer.Spec{
				tokenizer.NewSpec(`^\s+`, skipType),
				tokenizer.NewSpec(`^\(`, contextStartType),
				tokenizer.NewSpec(`^\)`, contextEndType),
				tokenizer.NewSpec(`^(\+|-|\*|\/|%)`, arithmeticOperationType),
				tokenizer.NewSpec(`^(==|!=|<=?|>=?)`, comparisonOperationType),
				tokenizer.NewSpec(`^(&&|\|\|)`, logicalOperationType),
				tokenizer.NewSpec(`^\d+(\.\d+)?`, numberType),
				tokenizer.NewSpec(`^(true|false)`, boolType),
				tokenizer.NewSpec(`^("(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*')`, textType),
			}),
	}
}

// eat return a token with expected type.
func (p *parser) eat(tokenType tokenizer.Type) (*tokenizer.Token, error) {
	token := p.lookahead1
	p.lookahead1 = p.lookahead2

	if token == nil {
		return nil, errs.NewErrUnexpectedInputEnd(tokenType.String())
	}

	if token.Type != tokenType {
		return nil, errs.NewErrorAtPosition(
			errs.NewErrUnexpectedTokenType(token.Type.String(), tokenType.String()),
			p.tokenizer.GetCursorPosition(),
		)
	}

	var err error
	p.lookahead2, err = p.tokenizer.GetNextToken()

	return token, err
}

func (p *parser) Parse() Result {
	var err error

	p.lookahead1, err = p.tokenizer.GetNextToken()
	if err != nil {
		return Result{
			Error: err,
		}
	}
	p.lookahead2, err = p.tokenizer.GetNextToken()
	if err != nil {
		return Result{
			Error: err,
		}
	}

	value, err := p.expression(nil)

	return Result{
		Value: value,
		Error: err,
	}
}

func (p *parser) expression(leftValue interface{}) (interface{}, error) {
	if p.lookahead1 == nil {
		return nil, errs.NewErrorAtPosition(
			errs.NewErrUnexpectedInputEnd("expression"),
			p.tokenizer.GetCursorPosition())
	}

	if p.lookahead1.Type == contextStartType {
		return p.contextExpression()
	}

	var err error
	if leftValue == nil {
		leftValue, err = p.literal()
		if p.lookahead1 == nil || err != nil {
			return leftValue, err
		}
	}

	if p.lookahead1.Type == arithmeticOperationType {
		return p.arithmeticOperation(leftValue)
	} else if p.lookahead1.Type == logicalOperationType {
		return p.logicalOperation(leftValue)
	} else if p.lookahead1.Type == comparisonOperationType {
		return p.comparisonOperation(leftValue)
	} else if p.lookahead1.Type == contextEndType {
		return leftValue, nil
	}

	return nil, errs.NewErrorAtPosition(
		errs.NewErrUnexpectedTokenType(p.lookahead1.Type.String(), "any"),
		p.tokenizer.GetCursorPosition())
}

func (p *parser) contextExpression() (interface{}, error) {
	_, err := p.eat(contextStartType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	value, err := p.expression(nil)
	if err != nil {
		return nil, err
	}

	_, err = p.eat(contextEndType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	if p.lookahead1 != nil {
		return p.expression(value)
	}

	return value, nil
}

func (p *parser) arithmeticOperation(leftValue interface{}) (interface{}, error) {
	token, err := p.eat(arithmeticOperationType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	if token.Value == "*" || token.Value == "/" || token.Value == "%" {
		rightValue, err := p.literal()
		if err != nil {
			return nil, err
		}

		rightValueConverted := convertFloat(rightValue)
		switch token.Value {
		case "*":
			rightValue = convertFloat(leftValue) * rightValueConverted
		case "/":
			if rightValueConverted == 0 {
				return nil, errs.NewErrorAtPosition(
					errs.ErrDivisionByZero,
					p.tokenizer.GetCursorPosition())
			}
			rightValue = convertFloat(leftValue) / rightValueConverted
		case "%":
			if rightValueConverted == 0 {
				return nil, errs.NewErrorAtPosition(
					errs.ErrDivisionByZero,
					p.tokenizer.GetCursorPosition())
			}
			rightValue = int64(math.Round(convertFloat(leftValue))) % int64(math.Round(rightValueConverted))
		}

		if p.lookahead1 == nil {
			return rightValue, nil
		}

		return p.expression(rightValue)
	}

	var rightValue interface{}
	if p.lookahead2 == nil || p.lookahead2.Type != arithmeticOperationType {
		rightValue, err = p.literal()
		if err != nil {
			return nil, err
		}
	} else {
		rightValue, err = p.expression(nil)
		if err != nil {
			return nil, err
		}
	}

	if token.Value == "+" {
		leftValue = convertFloat(leftValue) + convertFloat(rightValue)
	} else if token.Value == "-" {
		leftValue = convertFloat(leftValue) - convertFloat(rightValue)
	}

	if p.lookahead1 != nil {
		return p.expression(leftValue)
	} else {
		return leftValue, nil
	}
}

func (p *parser) logicalOperation(leftValue interface{}) (interface{}, error) {
	token, err := p.eat(logicalOperationType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	rightValue, err := p.expression(nil)
	if err != nil {
		return nil, err
	}

	switch token.Value {
	case "&&":
		return convertBool(leftValue) && convertBool(rightValue), nil
	case "||":
		return convertBool(leftValue) || convertBool(rightValue), nil
	}

	return nil, errs.NewErrorAtPosition(
		errs.NewErrUnexpectedTokenType(token.Type.String(), "logical operation"),
		p.tokenizer.GetCursorPosition())
}

func (p *parser) comparisonOperation(leftValue interface{}) (interface{}, error) {
	token, err := p.eat(comparisonOperationType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	rightValue, err := p.expression(nil)
	if err != nil {
		return nil, err
	}

	switch token.Value {
	case "==":
		return leftValue == rightValue, nil
	case "!=":
		return leftValue != rightValue, nil
	case "<":
		return convertFloat(leftValue) < convertFloat(rightValue), nil
	case "<=":
		return convertFloat(leftValue) <= convertFloat(rightValue), nil
	case ">":
		return convertFloat(leftValue) > convertFloat(rightValue), nil
	case ">=":
		return convertFloat(leftValue) >= convertFloat(rightValue), nil
	}

	return nil, errs.NewErrorAtPosition(
		errs.NewErrUnexpectedTokenType(token.Type.String(), "comparison operation"),
		p.tokenizer.GetCursorPosition())
}

func (p *parser) literal() (interface{}, error) {
	if p.lookahead1 == nil {
		return nil, errs.NewErrorAtPosition(
			errs.NewErrUnexpectedInputEnd("literal"),
			p.tokenizer.GetCursorPosition())
	}

	switch p.lookahead1.Type {
	case numberType:
		return p.number()
	case textType:
		return p.text()
	case boolType:
		return p.boolean()
	}

	return nil, errs.NewErrorAtPosition(
		errs.NewErrUnexpectedTokenType(p.lookahead1.Type.String(), "literal"),
		p.tokenizer.GetCursorPosition())
}

func (p *parser) number() (interface{}, error) {
	token, err := p.eat(numberType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	var value float64
	value, err = strconv.ParseFloat(token.Value, float64Size)
	if err != nil {
		return nil, errs.NewErrorAtPosition(
			fmt.Errorf("failed to parse float value: %w", err),
			p.tokenizer.GetCursorPosition()-len(token.Value))
	}

	return value, nil
}

func (p *parser) boolean() (interface{}, error) {
	token, err := p.eat(boolType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	return strings.ToLower(token.Value) == "true", nil
}

func (p *parser) text() (interface{}, error) {
	token, err := p.eat(textType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	return token.Value[1 : len(token.Value)-1], nil
}
