package atom

import (
	"go/token"
	"strings"
)

type Stream chan Atom

type Atom struct {
	Lit string
	token.Token
	token.Pos
}

func (a Atom) IsEOF() bool   { return a.Token == token.EOF }
func (a Atom) IsLegal() bool { return a.Token != token.ILLEGAL }

func (a Atom) String() string {
	if a.IsLiteral() {
		return a.Lit
	}
	return a.Token.String()
}

func Enquote(tok token.Token) string {
	if !tok.IsKeyword() && !tok.IsOperator() {
		return tok.String()
	}
	var b strings.Builder
	b.WriteRune('"')
	b.WriteString(tok.String())
	b.WriteRune('"')
	return b.String()
}
