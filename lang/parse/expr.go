package parse

import (
	"errors"
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

func (ex *Expr) parse(min oper.Level) (i item) {
	switch s := ex.Next(); {
	case s.IsIllegal(), s.IsEOF():
		return &stop{s}
	case s.IsKeyword():
		return &ctrl{s}
	case s.IsIdentifier(), s.IsLiteral():
		i = &term{s}
	default:
		if op, ok := oper.Default.Prefix(s.Token); ok {
			i = &rule{[]item{ex.parse(op.Level())}, op}
		} else {
			return &stop{s}
		}
	}
	for {
		op, ok := oper.Default.Get(ex.Next().Token,
			oper.UnaryLeft, oper.BinaryLeft, oper.BinaryRight)
		if !ok {
			return
		}
		l, r := op.Levels()
		if l.Compare(min) < 0 {
			ex.Streamer = ex.Undo(op.Symbol())
			return
		}
		a := []item{i}
		if r != oper.Unbound {
			a = append(a, ex.parse(r))
		}
		i = &rule{a, op}
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
