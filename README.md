# goeval
A GoLang library for simple expression evaluation.
```go
expr.Eval("true  && true") # true 
expr.Eval("2+3")   # 5 
expr.Eval("1+3>3") # true 
expr.Eval("('fo'!='bar')&&(2!=3)") # true 
```

[![GitHub release badge](https://badgen.net/github/release/StevenCyb/goeval/latest?label=Latest&logo=GitHub)](https://github.com/StevenCyb/goeval/releases/latest)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/StevenCyb/goeval/ci-test.yml?label=Tests&logo=GitHub)
![GitHub](https://img.shields.io/github/license/StevenCyb/goeval)

## Logical Operation
Supported logical operations are `&&` and `||`.
* On text, `"true"` or `'true'` (case insensitive) is `true`, else `false`
* On numbers greater zero is `true`, else `false`

## Comparison Operation
Supported comparisons are `==`, `!=`, `<`, `>`, `<=` or `>=`.
While equal and not equal directly uses the value. Greater(equal) and smaller(equal) will behave different:
* On boolean a `0` and `1` is used (`true>0` = `true`, `false>0` = `false`).
* On text the length is used (`"hello">"hey"` = `true`).

## Arithmetic Operation
Supported arithmetics are `+`, `-`, `*`, `/` and `%`.
* On boolean a `0` and `1` is used (`true+0` = `1`, `false+0` = `0`).
* On text the length is used (`"foo"+"bar"` = `6`).

## Context
Context can be used to group a part of the expression to prioritize the evaluation.