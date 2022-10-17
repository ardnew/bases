package op

import (
	"go/token"

	"github.com/ardnew/bases/io/syntax/parse/state"
)

// Operator provides an abstraction for any type of operator in prefix, postfix,
// infix, or other complex expressions.
type Operator struct {
	state.Parser
	L, R Level
	token.Token
}

func (p Operator) String() string { return p.Token.String() }

// Add adds an arbitrary number of given [token.Token] keys mapped to Operator
// values constructed with the given precedence, associativity, and
// [parse.Parser] state function.
func (m Table) Add(
	precedence int, assoc Assoc, par state.Parser, tok ...token.Token,
) {
	lhs, rhs := assoc.Level(precedence)
	for _, t := range tok {
		m[t] = Operator{Parser: par, L: lhs, R: rhs, Token: t}
	}
}

type Table map[token.Token]Operator

var Prefix, Postfix, Infix Table

func Init() {
	Infix.Add(22, BinaryRight, nil,
		token.PERIOD)

	Postfix.Add(21, UnaryLeft, nil,
		token.INC, token.DEC)

	Prefix.Add(20, UnaryRight, nil,
		token.INC, token.DEC)
	Prefix.Add(19, UnaryRight, nil,
		token.ADD, token.SUB)
	Prefix.Add(18, UnaryRight, nil,
		token.NOT, token.TILDE)

	Infix.Add(17, BinaryLeft, nil,
		token.MUL, token.QUO, token.REM)
	Infix.Add(16, BinaryLeft, nil,
		token.ADD, token.SUB)
	Infix.Add(15, BinaryLeft, nil,
		token.SHL, token.SHR)
	Infix.Add(14, BinaryLeft, nil,
		token.LSS, token.GTR, token.LEQ, token.GEQ)
	Infix.Add(13, BinaryLeft, nil,
		token.EQL, token.NEQ)
	Infix.Add(12, BinaryLeft, nil,
		token.AND)
	Infix.Add(11, BinaryLeft, nil,
		token.AND_NOT)
	Infix.Add(10, BinaryLeft, nil,
		token.XOR)
	Infix.Add(9, BinaryLeft, nil,
		token.OR)
	Infix.Add(8, BinaryLeft, nil,
		token.LAND)
	Infix.Add(7, BinaryLeft, nil,
		token.LOR)
	Infix.Add(6, BinaryRight, nil,
		token.DEFINE, token.ASSIGN)
	Infix.Add(5, BinaryRight, nil,
		token.ADD_ASSIGN, token.SUB_ASSIGN)
	Infix.Add(4, BinaryRight, nil,
		token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN)
	Infix.Add(3, BinaryRight, nil,
		token.SHL_ASSIGN, token.SHR_ASSIGN)
	Infix.Add(2, BinaryRight, nil,
		token.AND_ASSIGN, token.AND_NOT_ASSIGN, token.XOR_ASSIGN, token.OR_ASSIGN)
	Infix.Add(1, BinaryLeft, nil,
		token.COMMA, token.SEMICOLON)

	// Unhandled recognizable tokens:
	//   ARROW    // <-
	//   ELLIPSIS // ...
	//   COLON    // :

	//   LPAREN   // (
	//   LBRACK   // [
	//   LBRACE   // {
	//
	//   RPAREN   // )
	//   RBRACK   // ]
	//   RBRACE   // }
}
