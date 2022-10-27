package eval

import (
	"github.com/ardnew/bases/lang/lex/scan"
)

type eval struct {
	*scan.Scan
}

func EvalString(_ string) error {
	e := eval{}
	return e.Scan
}

//func (e *eval) Parse() error {
//	for e.Parser != nil {
//		e.Parser = e.Parser.Parse(e)
//	}
//	return e
//}
