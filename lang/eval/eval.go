package eval

import (
	"fmt"
	"go/token"
	"io"

	"github.com/ardnew/bases/lang/lex"
	"github.com/ardnew/bases/lang/lex/scan"
	"github.com/ardnew/bases/lang/parse"
	"github.com/ardnew/bases/lang/parse/expr"
	"github.com/ardnew/bases/lang/parse/oper"
)

type Eval struct {
	lex.Lexer
	parse.Parser
}

func New() *Eval {
	return &Eval{
		Lexer:  scan.New(),
		Parser: expr.New(),
	}
}

func (e *Eval) ReadFrom(r io.Reader) (n int64, err error) {
	return e.Parse(r)
}

func EvalString(s string) error {
	e := eval{Scan: scan.New().Init([]byte(s))}
	go e.Until(token.EOF, token.ILLEGAL)
	return e.Parse()
}

func (e *eval) Parse() error {
	fmt.Printf("Expr = %s\n", expr.Climb(e, oper.Unbound))
	return e
}
