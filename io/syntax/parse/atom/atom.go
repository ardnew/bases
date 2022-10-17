package atom

import "go/token"

type Stream chan Atom

type Atom struct {
	Lit string
	token.Token
	token.Pos
}

func (t Atom) String() string {
	if t.IsLiteral() {
		return t.Lit
	}
	return t.Token.String()
}
