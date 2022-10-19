package eval

import (
	"github.com/ardnew/bases/io/syntax/parse"
	"github.com/ardnew/bases/io/syntax/parse/expr"
	"github.com/ardnew/bases/io/syntax/parse/lex/scan"
)

type eval struct {
	parse.Parser
	*scan.Scan
}

func EvalString(s string) error {
	e := eval{
		Parser: expr.Parse,
		Scan:   scan.NewScan().Init([]byte(s)),
	}
	return e.Parse()
}

func (e *eval) Parse() error {
	return e
}
