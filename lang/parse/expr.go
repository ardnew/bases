package parse

import (
	"errors"
	"go/token"
	"io"
	"strings"

	"github.com/ardnew/bases/lang/parse/oper"
	"github.com/ardnew/bases/lang/parse/sym"
)

var ErrInvalidExpr = errors.New("invalid expression")

// var elog = log.New(log.DefaultWriter, "%f:%n ┆ %s")

// Expr represents an expression tree.
type Expr struct {
	item
	sym.Streamer
}

func New() *Expr { return &Expr{} }

func (e *Expr) Parse(r io.Reader) (n int64, err error) {
	return
}

func (e *Expr) ParseBuffer(b []byte) (n int64, err error) {
	logf("ParseBuffer(%+v): %q", b, string(b))
	e.Streamer = sym.Stream(b, sym.IsEOF, sym.IsIllegal)
	e.item = e.parse(oper.Unbound)
	return
}

func (e *Expr) ParseString(s string) (n int64, err error) {
	return e.ParseBuffer([]byte(s))
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

func (ex *Expr) parse(min oper.Level) item {
	symbol := ex.Next()
	switch it := wrap(symbol).(type) {
	case *stop:
		return it
	case *term:
		return it
	case *ctrl:
		return it
	case *rule:
		var pre bool
		switch it.Operator, pre = oper.Default.Prefix(symbol.Token); {
		case pre:
			_, br := it.Level()
			it.arg = append(it.arg, ex.parse(br))
		}
		return it
	}
	return nil
}

func (ex *Expr) Climb(depth int, min oper.Level) (it item) {
	s := ex.Next()
	l := wrap(s)
	logf("%*s%s -> %T:", depth*2, "", s, l)
	switch e := l.(type) {
	case *stop:

	case *term, *ctrl:
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
				if oper.Compare(bl, min) < 0 {
					break
				}
				ex.Next()
				l = newRule(os, l)
				continue
			} else if op, ok := oper.Default.Infix(os.Token); ok {
				bl, rl := op.Level()
				if oper.Compare(bl, min) < 0 {
					break
				}
				ex.Next()
				r := ex.Climb(depth+1, rl)
				l = newRule(os, l, r)
				continue
			} else {
				// Error!
				logf("Unexpected token: %T: %q [% x]", os, os, os)
				ex.Next()
				break
			}
		}
	}
	logf("%*s= %s", depth*2, "", l)
	return l
}

func newRule(s sym.Symbol, it ...item) *rule {
	return &rule{
		Operator: oper.Wrap(s),
		arg:      append(make([]item, 0, oper.MaxArity), it...),
	}
}

func (r *rule) String() string {
	var b strings.Builder
	b.WriteRune('(')
	b.WriteString(r.Operator.String())
	for _, a := range r.arg {
		b.WriteRune(' ')
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
