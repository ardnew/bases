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

func (a Atom) QQ() string { return a.enquote('"', '"') }
func (a Atom) Q() string  { return a.enquote('\'', '\'') }
func (a Atom) R() string  { return a.enquote('`', '`') }

func (a Atom) enquote(lhs, rhs rune) string {
	if !a.IsLiteral() || len(a.Lit) == 0 {
		return enquote(a.Token, lhs, rhs)
	}
	// Surround literals with a pair of braces
	lhl, rhl := '{', '}'
	if lhs == lhl && rhs == rhl {
		// Fallback if our braces are the same as the quotes
		lhl, rhl = '<', '>'
	}
	var b strings.Builder
	b.WriteRune(lhl)
	b.WriteRune(lhs)
	switch a.Token {
	case token.CHAR, token.STRING:
		lit := []rune(a.Lit)
		// Remove any existing quotes surrounding the literal
		if n := len(lit); n > 1 && a.isQuote(lit[0]) && lit[0] == lit[n-1] {
			lit = lit[1 : len(lit)-1]
		}
		b.WriteString(string(lit))

	case token.IDENT, token.INT, token.FLOAT, token.IMAG:
		fallthrough
	default:
		b.WriteString(a.Lit)
	}
	b.WriteRune(rhs)
	b.WriteRune(rhl)
	return b.String()
}

func (a Atom) isQuote(u rune) bool {
	switch a.Token {
	case token.CHAR:
		return u == '\''
	case token.STRING:
		return u == '"' || u == '`'
	default:
		return false
	}
}

func QQ(tok token.Token) string { return Atom{Token: tok}.QQ() }
func Q(tok token.Token) string  { return Atom{Token: tok}.Q() }
func R(tok token.Token) string  { return Atom{Token: tok}.R() }

func enquote(tok token.Token, lhs, rhs rune) string {
	var b strings.Builder
	switch tok {
	// Non-representable tokens enclosed in angles instead of quotes
	case token.ILLEGAL, token.EOF:
		b.WriteRune('<')
		b.WriteString(tok.String())
		b.WriteRune('>')

	case token.COMMENT:
		b.WriteString("/*")
		b.WriteString(tok.String())
		b.WriteString("*/")

	// Literal primitives
	case token.IDENT,
		token.INT, token.FLOAT, token.IMAG,
		token.CHAR, token.STRING:
		fallthrough

	// Keywords, operators, and unrecognized tokens
	default:
		b.WriteRune(lhs)
		b.WriteString(tok.String())
		b.WriteRune(rhs)
	}
	return b.String()
}
