package eval

import (
	"fmt"
	"go/token"

	"github.com/ardnew/bases/lang/lex/scan"
	"github.com/ardnew/bases/lang/parse"
	"github.com/ardnew/bases/lang/parse/op"
)

type eval struct {
	*scan.Scan
}

func EvalString(s string) error {
	e := eval{Scan: scan.New().Init([]byte(s))}
	go e.Until(token.EOF, token.ILLEGAL)
	return e.Parse()
}

func (e *eval) Parse() error {
	fmt.Printf("Expr = %s\n", parse.Climb(e, op.Unbound))
	return e
}
