package eval

import (
	"github.com/ardnew/bases/io/syntax/parse"
	"github.com/ardnew/bases/io/syntax/parse/expr"
	"github.com/ardnew/bases/io/syntax/parse/lex/scan"
)

type eval struct {
	parse.State
	*scan.Scan
}

func EvalString(s string) error {
	e := eval{
		State: expr.Parse,
		Scan:  scan.NewScan().Init([]byte(s)),
	}
	return e.Parse()
}

func (e *eval) Parse() error {
	for e.State != nil {
		e.State = e.State(e)
	}
	return e
}
