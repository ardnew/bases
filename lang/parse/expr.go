package parse

import (
	"errors"
	"go/token"
	"io"
	"strings"

	"github.com/ardnew/bases/lang/parse/oper"
	"github.com/ardnew/bases/lang/parse/sym"
	"github.com/ardnew/bases/log"
)

var ErrInvalidExpr = errors.New("invalid expression")

// var elog = log.New(log.DefaultWriter, "%f:%n ┆ %s")

// Expr represents an expression tree.
type Expr struct {
	item
	sym.Streamer
	log *log.Log
}

func New() *Expr {
	p := &Expr{}
	if w, f, e := log.LookupEnv("EXPR"); e == nil {
		p.log = log.New(w, f)
		p.log.SetCallerOffset(2)
	}
	return p
}

func (e *Expr) logf(format string, v ...interface{}) {
	if e.log != nil {
		e.log.Printf(format, v...)
	}
}

func (e *Expr) Parse(r io.Reader) (n int64, err error) {
	return
}

func (e *Expr) ParseBuffer(b []byte) (n int64, err error) {
	e.logf("ParseBuffer(%+v): %q", b, string(b))
	e.Streamer = sym.Stream(b, sym.IsEOF, sym.IsIllegal)
	e.item = e.Climb(0, oper.Unbound)
	return
}

func (e *Expr) ParseString(s string) (n int64, err error) {
	return
}

func (e *Expr) Err() error {
	if e == nil || e.item == nil {
		return ErrInvalidExpr
	}
	return nil
}

type (
	item interface{ String() string }
	rule struct {
		arg []item
		oper.Operator
	}
	stop struct{ sym.Symbol }
	term struct{ sym.Symbol }
	ctrl struct{ sym.Symbol }
)

func wrap(s sym.Symbol) item {
	switch {
	case s.IsIllegal(), s.IsEOF():
		return &stop{s}
	case s.IsIdentifier(), s.IsLiteral():
		return &term{s}
	case s.IsKeyword():
		return &ctrl{s}
	default:
		return newRule(s)
	}
}

func (ex *Expr) Climb(depth int, min oper.Level) (it item) {
	s := ex.Next()
	l := wrap(s)
	ex.logf("%*s%s -> %T:", depth*2, "", s, l)
	switch e := l.(type) {
	case *stop, *term, *ctrl:
	case *rule:
		var prefix bool
		switch e.Operator, prefix = oper.Default.Prefix(s.Token); {
		case e.Spells(token.LPAREN):
			l = ex.Climb(depth+1, oper.Unbound)
			if t := ex.Next(); !t.Is(sym.Operator(token.RPAREN)) {
				ex.Streamer = ex.Undo(t)
			}
		case prefix:
			_, br := e.Level()
			e.arg = append(e.arg, ex.Climb(depth+1, br))
		default:
		}
	}

	for {
		if os := ex.Peek(); os.IsEOF() {
			break
		} else {
			if op, ok := oper.Default.Postfix(os.Token); ok {
				bl, _ := op.Level()
				var bm int
				if min != oper.Unbound {
					bm = min.Int()
				}
				if bl.Int() < bm {
					break
				}
				ex.Next()
				l = newRule(os, l)
				continue
			} else {
				ex.Next()
				break
			}
		}
	}
	ex.logf("%*s= %s", depth*2, "", l)
	return l
}

func newRule(s sym.Symbol, it ...item) *rule {
	return &rule{
		Operator: oper.Wrap(s),
		arg:      append(make([]item, 0, oper.MaxArity), it...),
	}
}

// func (t *Rule) Parse(lexer lex.Lexer, level oper.Level) Expr {
// 	ls := lexer.Take()
// 	switch oper, prefix := oper.Default.Prefix(ls.Token); {
// 	case prefix && oper.Token() == token.LPAREN:
// 		if e := t.Parse(lexer, oper.Unbound); e == nil {
// 			// unexpected bytes or EOF
// 		}
// 		if !lexer.Check(sym.Token(token.RPAREN)) {
// 			// unclosed paren
// 		}

// 	case prefix:
// 		t.Op = oper
// 		if e := t.Parse(lexer, oper.Right()); e == nil {
// 			//
// 		}
// 		t.E = []Expr{}

// 		e = &Rule{Op: oper, E: []Expr{climb(lexer, oper.Right())}}
// 	default:
// 		e = &Term{lhs}
// 	}

// 	return false
// }

func (r *rule) String() string {
	var b strings.Builder
	b.WriteRune('(')
	b.WriteString(r.Operator.String())
	for _, a := range r.arg {
		b.WriteString(", ")
		b.WriteString(a.String())
	}
	b.WriteRune(')')
	return b.String()
}

func (s *stop) String() string {
	return s.Symbol.String()
}

func (t *term) String() string {
	return t.Symbol.String()
}

func (c *ctrl) String() string {
	return c.Symbol.String()
}
