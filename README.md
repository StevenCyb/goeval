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

