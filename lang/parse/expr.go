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

// Expr represents an expression tree.
type Expr struct {
	item
	sym.Streamer
	in chan sym.Symbol

	nb int
}

func New() *Expr { return &Expr{in: make(chan sym.Symbol)} }

func (e *Expr) Parse(r io.Reader) (n int64, err error) {
	return
}

func (e *Expr) ParseBuffer(b []byte) (n int64, err error) {
	e.Streamer = sym.Stream(b)
	e.Go(e.in, sym.IsEOF, sym.IsIllegal)
	e.nb = 0
	e.item = e.Climb(oper.Unbound)
	return int64(e.nb), e.Err()
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

func (ex *Expr) Climb(min oper.Level) (it item) {
	s := <-ex.in
	defer func(nb int) {
		if ex.nb < nb {
			ex.nb = nb
		}
	}(int(s.Pos) + len(s.String()))
	l := wrap(s)
	switch e := l.(type) {
	case *stop, *term, *ctrl:
	case *rule:
		var prefix bool
		switch e.Operator, prefix = oper.Default.Prefix(s.Token); {
		case e.Spells(token.LPAREN):
			l = ex.Climb(oper.Unbound)
			if t := <-ex.in; !t.Is(sym.Operator(token.RPAREN)) {
				ex.Streamer = ex.Undo(t)
			}
		case prefix:
			_, br := e.Level()
			e.arg = append(e.arg, ex.Climb(br))
		default:
		}
	}

	// for {
	// 	if os := lexer.Look(); os.IsEOF() {
	// 		break
	// 	} else {
	// 		if op, ok := oper.Default.Postfix(os.Token); ok {
	// 			bl, _ := op.Level()
	// 			if bl.Int() < min.Int() {
	// 				break
	// 			}
	// 			lexer.Take()
	// 			l = newRule(os, l)
	// 			continue
	// 		} else {
	// 			lexer.Take()
	// 			break
	// 		}
	// 	}
	// }

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
