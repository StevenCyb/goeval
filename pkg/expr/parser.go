package expr

import (
	"fmt"
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

// LL(2) parser for the following grammar:
/*
 * <EXPRESSION>            ::= <SKIP>* <NUMBER_EXPANDED> | <SKIP>* <TEXT_EXPANDED> | <SKIP>* <LOGICAL_EXPRESSION> | <SKIP>* <COMPARISON_EXPRESSION>
 *
 * <BOOL_EXPANDED>         ::= <BOOL> <SKIP>* | <BOOL_EXPANDED> <LOGICAL_OPERATION> <SKIP>* <BOOL_EXPANDED> | <BOOL_EXPANDED> <LOGICAL_OPERATION> <BOOL_EXPANDED>
 * <TEXT_EXPANDED>         ::= <TEXT> <SKIP>* | <TEXT> <SKIP>* "+" <SKIP>* <TEXT_EXPANDED>
 * <NUMBER_EXPANDED>       ::= <NUMBER> <SKIP>* | <NUMBER> <SKIP>* <ARITHMETIC_OPERATION> <SKIP>* <NUMBER> <SKIP>* | <NUMBER> <SKIP>* <ARITHMETIC_OPERATION> <NUMBER_EXPANDED>
 *
 * <COMPARISON_EXPRESSION> ::= <EXPRESSION> <COMPARISON_OPERATION> <EXPRESSION> | <EXPRESSION> <COMPARISON_OPERATION> <COMPARISON_EXPRESSION>

# TODO CONTEXT_EXPRESSION

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

// Create a new LL(1) parser for the given expression.
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
				tokenizer.NewSpec(`^("[^"]*"|'[^"]*')`, textType),
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

	value, err := p.expression()

	return Result{
		Value: value,
		Error: err,
	}
}

func (p *parser) expression() (interface{}, error) {
	if p.lookahead1 == nil {
		return nil, errs.NewErrorAtPosition(
			errs.NewErrUnexpectedInputEnd("expression"),
			p.tokenizer.GetCursorPosition())
	}

	switch p.lookahead1.Type {
	case numberType:
		return p.numberExpanded()
	case textType:
		return p.textExpanded()
	case boolType:
		return p.boolExpanded()
	}

	return nil, errs.NewErrorAtPosition(
		errs.NewErrUnexpectedTokenType(p.lookahead1.Type.String(), "expression"),
		p.tokenizer.GetCursorPosition()-len(p.lookahead1.Value))
}

func (*parser) comparisonExpression() (interface{}, error) {

	return nil, nil
}

func (p *parser) numberExpanded() (interface{}, error) {
	firstValue, err := p.number()
	if err != nil {
		return nil, err
	}

	if p.lookahead1 != nil && p.lookahead1.Type == arithmeticOperationType {
		operation, err := p.eat(arithmeticOperationType)
		if err != nil {
			return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
		}

		if operation.Value == "*" || operation.Value == "/" || operation.Value == "%" {
			secondValue, err := p.number()
			if err != nil {
				return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
			}

			switch operation.Value {
			case "*":
				firstValue = firstValue.(float64) * secondValue.(float64)
			case "/":
				firstValue = firstValue.(float64) / secondValue.(float64)
			case "%":
				firstValue = int64(firstValue.(float64)) % int64(secondValue.(float64))
			}

			if p.lookahead1 == nil || p.lookahead1.Type != arithmeticOperationType {
				return firstValue, nil
			}

			operation, err = p.eat(arithmeticOperationType)
			if err != nil {
				return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
			}
		}

		secondValue, err := p.numberExpanded()
		if err != nil {
			return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
		}

		switch operation.Value {
		case "+":
			return firstValue.(float64) + secondValue.(float64), nil
		case "-":
			return firstValue.(float64) - secondValue.(float64), nil
		case "*":
			return firstValue.(float64) * secondValue.(float64), nil
		case "/":
			return firstValue.(float64) / secondValue.(float64), nil
		case "%":
			return int64(firstValue.(float64)) % int64(secondValue.(float64)), nil
		}
	}

	return firstValue, nil
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

func (p *parser) boolExpanded() (interface{}, error) {
	firstValue, err := p.boolean()
	if err != nil {
		return nil, err
	}

	if p.lookahead1 != nil && p.lookahead1.Type == logicalOperationType {
		operation, err := p.eat(logicalOperationType)
		if err != nil {
			return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
		}

		secondValue, err := p.boolExpanded()
		if err != nil {
			return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
		}

		switch operation.Value {
		case "&&":
			return firstValue.(bool) && secondValue.(bool), nil
		case "||":
			return firstValue.(bool) || secondValue.(bool), nil
		}
	}

	return firstValue, nil
}

func (p *parser) boolean() (interface{}, error) {
	token, err := p.eat(boolType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	return strings.ToLower(token.Value) == "true", nil
}

func (p *parser) textExpanded() (interface{}, error) {
	token, err := p.eat(textType)
	if err != nil {
		return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
	}

	token.Value = strings.ReplaceAll(token.Value, "\"", "")
	token.Value = strings.ReplaceAll(token.Value, "'", "")

	if p.lookahead1 != nil && p.lookahead1.Type == arithmeticOperationType && p.lookahead1.Value == "+" {
		_, err = p.eat(arithmeticOperationType)
		if err != nil {
			return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
		}

		secondValue, err := p.textExpanded()
		if err != nil {
			return nil, errs.NewErrorAtPosition(err, p.tokenizer.GetCursorPosition())
		}

		return fmt.Sprintf("%s%s", token.Value, secondValue), nil
	}

	return token.Value, nil
}
