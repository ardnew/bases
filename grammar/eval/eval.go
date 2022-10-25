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
		Parser: expr.SExpr(""),
		Scan:   scan.NewScan().Init([]byte(s)),
	}
	for e.Parser != nil {
		e.Parser = e.Parse(e)
	}
	return e.Scan
}

//func (e *eval) Parse() error {
//	for e.Parser != nil {
//		e.Parser = e.Parser.Parse(e)
//	}
//	return e
//}
