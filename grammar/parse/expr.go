package parse

import (
	"go/token"
	"strings"

	"github.com/ardnew/bases/grammar/lex"
	"github.com/ardnew/bases/grammar/lex/sym"
	"github.com/ardnew/bases/grammar/parse/op"
)

type E interface {
	String() string
}

// func climb(lexer lex.Lexer, level op.Level) (e E) {
// 	lhs := lexer.Take()
// 	switch oper, prefix := op.Default.Prefix(lhs.Token); {
// 	case prefix && oper.Token() == token.LPAREN:
// 		e = climb(lexer, op.Unbound)
// 		if !lexer.Check(sym.Token(token.RPAREN)) {
// 			// unclosed paren
// 		}
// 	case prefix:
// 		e = &Expr{Op: oper, E: []E{climb(lexer, oper.Right())}}
// 	default:
// 		e = &Term{lhs}
// 	}

// 	for {
// 		if os := lexer.Look(); os.IsEOF() {
// 		} else {
// 			if oper, ok := op.Default.Postfix(os.Token); ok {
// 				// e =
// 			} else {
// 			}
// 		}
// 	}

// 	return e
// }

type Expr struct {
	E  []E
	Op op.Operator
}

type Term struct {
	sym.Symbol
}

func (t *Expr) Parse(lexer lex.Lexer, level op.Level) E {
	ls := lexer.Take()
	switch oper, prefix := op.Default.Prefix(ls.Token); {
	case prefix && oper.Token() == token.LPAREN:
		if e := t.Parse(lexer, op.Unbound); e == nil {
			// unexpected bytes or EOF
		}
		if !lexer.Check(sym.Token(token.RPAREN)) {
			// unclosed paren
		}

	case prefix:
		t.Op = oper
		if e := t.Parse(lexer, oper.Right()); e == nil {
			//
		}
		t.E = []E{}

		e = &Expr{Op: oper, E: []E{climb(lexer, oper.Right())}}
	default:
		e = &Term{lhs}
	}

	return false
}

func (t *Term) Parse(lexer lex.Lexer, level op.Level) E {
	switch s := lexer.Look(); {
	case s.IsTerminal():
		*t = Term{lexer.Take()}
		return t
	case s.IsEOF() || !s.IsLegal():
		*t = Term{lexer.Take()}
		return nil
	default:
		return nil
	}
}

func (s *Expr) String() string {
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

func (t *Term) String() string {
	return t.Symbol.String()
}
