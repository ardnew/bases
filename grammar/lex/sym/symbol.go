package sym

import (
	"go/token"
)

type Stream chan Symbol

type Symbol struct {
	Lit string
	token.Token
	token.Pos
}

func Token(tok token.Token) Symbol {
	return Symbol{Token: tok}
}

func Literal(tok token.Token, lit string) Symbol {
	return Symbol{Lit: lit, Token: tok}
}

func (a Symbol) IsEOF() bool   { return a.Token == token.EOF }
func (a Symbol) IsLegal() bool { return a.Token != token.ILLEGAL }

func (a Symbol) IsFollowed() bool {
	return a.IsLegal() && !a.IsEOF()
}

func (a Symbol) IsTerminal() bool {
	return a.IsLiteral() || a.IsKeyword() ||
		token.IsIdentifier(a.Lit)
}

func (a Symbol) Matches(t Symbol) bool {
	return a.Token == t.Token && (!a.IsLiteral() || a.Lit == t.Lit)
}
