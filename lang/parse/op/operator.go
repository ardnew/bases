package op

import (
	"go/token"

	"github.com/ardnew/bases/lang/lex/sym"
)

// Operator provides an abstraction for any type of operator in prefix, postfix,
// infix, or other complex expressions.
type Operator struct {
	lbl, rbl Level
	sym      sym.Symbol
	prc      int
	ass      Assoc
}

func (o Operator) Level() (Level, Level)   { return o.lbl, o.rbl }
func (o Operator) Symbol() sym.Symbol      { return o.sym }
func (o Operator) Is(tok token.Token) bool { return o.sym.Token == tok }
func (o Operator) String() string          { return o.sym.String() }

func Wrap(sym sym.Symbol) Operator {
	return Operator{sym: sym}
}

type operators [maxOperators]Operator

func (p *operators) get(tok token.Token) (op Operator, ok bool) {
	if 0 <= int(tok) && int(tok) < maxOperators {
		op, ok = p[tok], (p[tok].lbl != Unbound) || (p[tok].rbl != Unbound)
	}
	return
}

func (p *operators) add(prc int, ass Assoc, tok ...token.Token) {
	lhs, rhs := ass.Level(prc)
	for _, t := range tok {
		p[t] = Operator{
			lbl: lhs, rbl: rhs,
			sym: sym.Operator(t),
			prc: prc, ass: ass,
		}
	}
}
