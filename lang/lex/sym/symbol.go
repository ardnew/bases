package sym

import (
	"go/token"
)

// Symbol is a [token.Token] represented by a string at position [token.Pos].
type Symbol struct {
	Lit string
	token.Token
	token.Pos
}

// Illegal returns a Symbol with [token.ILLEGAL].
func Illegal() Symbol {
	return Symbol{Token: token.ILLEGAL}
}

// EOF returns a Symbol with [token.EOF].
func EOF() Symbol {
	return Symbol{Token: token.EOF}
}

// Operator returns a Symbol with the given token and an empty string literal.
func Operator(tok token.Token) Symbol {
	return Symbol{Token: tok}
}

// Keyword returns a Symbol with the given token and an empty string literal.
func Keyword(tok token.Token) Symbol {
	return Symbol{Token: tok}
}

// Identifier returns a Symbol with [token.IDENT] and the given string literal.
func Identifier(ident string) Symbol {
	return Symbol{Lit: ident, Token: token.IDENT}
}

// Literal returns a Symbol with the given token and string literal for one of
// the five basic types:
// [token.INT], [token.FLOAT], [token.IMAG], [token.CHAR], or [token.STRING].
func Literal(tok token.Token, lit string) Symbol {
	return Symbol{Lit: lit, Token: tok}
}

// IsIllegal returns true if the token is equal to [token.ILLEGAL];
// it returns false otherwise.
func (s Symbol) IsIllegal() bool {
	return s.Token == token.ILLEGAL
}

// IsEOF returns true if the token is equal to [token.EOF];
// it returns false otherwise.
func (s Symbol) IsEOF() bool {
	return s.Token == token.EOF
}

// IsIdentifier returns true if the token is a valid Go identifier;
// it returns false otherwise.
//
// The token is a valid identifier if it is equal to [token.IDENT] and
// token.IsIdentifier(s.Lit) returns true (see [go/token.IsIdentifier]).
func (s Symbol) IsIdentifier() bool {
	return s.Token == token.IDENT && token.IsIdentifier(s.Lit)
}

// IsLiteral returns true if the token is one of the five basic types:
// [token.INT], [token.FLOAT], [token.IMAG], [token.CHAR], or [token.STRING];
// it returns false otherwise.
//
// Unlike [go/token.Token.IsLiteral], identifiers are not considered literals.
// In other words, s.IsIdentifier() and s.IsLiteral() are mutually exclusive.
func (s Symbol) IsLiteral() bool {
	return s.Token != token.IDENT && s.Token.IsLiteral()
}

// Is returns true if s and t represent the same token without respect to
// location; it returns false otherwise.
//
// In particular, a.Is(b) does not imply a.Pos == b.Pos.
func (s Symbol) Is(t Symbol) bool {
	return s.Token == t.Token && (!s.Token.IsLiteral() || s.Lit == t.Lit)
}
