package parse

import (
	"go/token"
	"strings"

	"github.com/ardnew/bases/lang/lex"
	"github.com/ardnew/bases/lang/lex/sym"
	"github.com/ardnew/bases/lang/parse/op"
)

type Expr interface {
	String() string
}

type (
	Rule struct {
		E  []Expr
		Op op.Operator
	}
	Stop struct{ sym.Symbol }
	Term struct{ sym.Symbol }
	Ctrl struct{ sym.Symbol }
)

func wrap(s sym.Symbol) Expr {
	switch {
	case s.IsIllegal(), s.IsEOF():
		return &Stop{s}
	case s.IsIdentifier(), s.IsLiteral():
		return &Term{s}
	case s.IsKeyword():
		return &Ctrl{s}
	default:
		return newRule(s)
	}
}

func Climb(lexer lex.Lexer, min op.Level) (expr Expr) {
	s := lexer.Take()
	l := wrap(s)
	switch e := l.(type) {
	case *Stop, *Term, *Ctrl:
	case *Rule:
		var prefix bool
		switch e.Op, prefix = op.Default.Prefix(s.Token); {
		case e.Op.Is(token.LPAREN):
			l = Climb(lexer, op.Unbound)
			if !lexer.Check(sym.Operator(token.RPAREN)) {
				// unclosed paren
			}
		case prefix:
			_, br := e.Op.Level()
			e.E = append(e.E, Climb(lexer, br))
		default:
		}

	}

	for {
		if os := lexer.Look(); os.IsEOF() {
			break
		} else {
			if oper, ok := op.Default.Postfix(os.Token); ok {
				bl, _ := oper.Level()
				if bl.Int() < min.Int() {
					break
				}
				lexer.Take()
				l = newRule(os, l)
				continue
			} else {
				lexer.Take()
				break
			}
		}
	}

	return l
}

func newRule(s sym.Symbol, e ...Expr) *Rule {
	return &Rule{append(make([]Expr, 0, op.MaxArity), e...), op.Wrap(s)}
}

// func (t *Rule) Parse(lexer lex.Lexer, level op.Level) Expr {
// 	ls := lexer.Take()
// 	switch oper, prefix := op.Default.Prefix(ls.Token); {
// 	case prefix && oper.Token() == token.LPAREN:
// 		if e := t.Parse(lexer, op.Unbound); e == nil {
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

// func (t *Term) Parse(lexer lex.Lexer, level op.Level) Expr {
// 	switch s := lexer.Look(); {
// 	case s.IsTerminal():
// 		*t = Term{lexer.Take()}
// 		return t
// 	case s.IsEOF() || !s.IsLegal():
// 		*t = Term{lexer.Take()}
// 		return nil
// 	default:
// 		return nil
// 	}
// }

func (s *Rule) String() string {
	var b strings.Builder
	b.WriteRune('(')
	b.WriteString(s.Op.String())
	for _, e := range s.E {
		b.WriteString(", ")
		b.WriteString(e.String())
	}
	b.WriteRune(')')
	return b.String()
}

func (s *Stop) String() string {
	return s.Symbol.String()
}

func (t *Term) String() string {
	return t.Symbol.String()
}

func (t *Ctrl) String() string {
	return t.Symbol.String()
}
