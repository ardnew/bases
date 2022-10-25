package op

import (
	"go/token"
)

// Operator provides an abstraction for any type of operator in prefix, postfix,
// infix, or other complex expressions.
type Operator struct {
	lbl, rbl Level
	prc      int
	ass      Assoc
	tok      token.Token
}

func (o Operator) Precedence() int         { return o.prc }
func (o Operator) Associativity() Assoc    { return o.ass }
func (o Operator) Order() (int, Assoc)     { return o.prc, o.ass }
func (o Operator) Binding() (Level, Level) { return o.lbl, o.rbl }
func (o Operator) Left() Level             { return o.lbl }
func (o Operator) Right() Level            { return o.rbl }
func (o Operator) Token() token.Token      { return o.tok }
func (o Operator) String() string          { return o.tok.String() }

type operators [maxOperators]Operator

func (p *operators) get(tok token.Token) (op Operator, ok bool) {
	if 0 <= int(tok) && int(tok) < maxOperators {
		op, ok = p[tok], (p[tok].lbl != Unbound) || (p[tok].rbl != Unbound)
	}
	return
}

func (p *operators) add(prec int, assoc Assoc, tok ...token.Token) {
	lhs, rhs := assoc.Level(prec)
	for _, t := range tok {
		p[t] = Operator{lbl: lhs, rbl: rhs, prc: prec, ass: assoc, tok: t}
	}
}
