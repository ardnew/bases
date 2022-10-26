package sym

import (
	"go/token"
)

// Stream represents a stream of input tokens.
type Stream chan Symbol

// Symbol is a [token.Token] represented by the string Lit at byte offset
// [token.Pos].
//
// When the location of a Symbol is not important, Pos will retain its zero
// value with the same semantics as described by [go/token.Pos].
type Symbol struct {
	Lit string
	token.Token
	token.Pos
}

// Illegal returns the ILLEGAL Symbol.
func Illegal() Symbol {
	return Symbol{Token: token.ILLEGAL}
}

// EOF returns the EOF Symbol.
func EOF() Symbol {
	return Symbol{Token: token.EOF}
}

// Operator returns an operator Symbol with the given token and an empty string
// literal.
func Operator(tok token.Token) Symbol {
	return Symbol{Token: tok}
}

// Keyword returns a keyword Symbol with the given token and an empty string
// literal.
func Keyword(tok token.Token) Symbol {
	return Symbol{Token: tok}
}

// Identifier returns an IDENT Symbol with the given string literal.
func Identifier(ident string) Symbol {
	return Symbol{Lit: ident, Token: token.IDENT}
}

// Literal returns a Symbol with the given token and string literal for one of
// the five basic types: INT, FLOAT, IMAG, CHAR, or STRING.
func Literal(tok token.Token, lit string) Symbol {
	return Symbol{Lit: lit, Token: tok}
}

// IsIllegal returns whether the token is equal to [token.ILLEGAL].
func (s Symbol) IsIllegal() bool {
	return s.Token == token.ILLEGAL
}

// IsEOF returns whether the token is equal to [token.EOF].
func (s Symbol) IsEOF() bool {
	return s.Token == token.EOF
}

// IsIdentifier returns whether the token is a valid Go identifier.
func (s Symbol) IsIdentifier() bool {
	return s.Token == token.IDENT && token.IsIdentifier(s.Lit)
}

// IsLiteral returns whether the token is a literal of one of the five basic
// types: INT, FLOAT, IMAG, CHAR, or STRING.
func (s Symbol) IsLiteral() bool {
	return s.Token != token.IDENT && s.Token.IsLiteral()
}

// Matches tests whether t and s represent the same token without respect to
// location, i.e., s.Pos may be unequal to t.Pos.
func (s Symbol) Matches(t Symbol) bool {
	return s.Token == t.Token && (!s.Token.IsLiteral() || s.Lit == t.Lit)
}
