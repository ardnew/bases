package op

import (
	"go/token"
)

func init() { Default.Reset() }

// Prefix calls method [Default.Prefix], which is the default prefix operator
// precedence table.
func Prefix(tok token.Token) (Operator, bool) { return Default.Prefix(tok) }

// Postfix calls method [Default.Postfix], which is the default postfix operator
// precedence table.
func Postfix(tok token.Token) (Operator, bool) { return Default.Postfix(tok) }

// Infix calls method [Default.Infix], which is the default infix operator
// precedence table.
func Infix(tok token.Token) (Operator, bool) { return Default.Infix(tok) }

// Operator provides an abstraction for any type of operator in prefix, postfix,
// infix, or other complex expressions.
type Operator struct {
	bll, blr Level
	prc      int
	ass      Assoc
	tok      token.Token
}

func (p Operator) String() string { return p.tok.String() }

func (p Operator) Precedence() int      { return p.prc }
func (p Operator) Associativity() Assoc { return p.ass }
func (p Operator) Order() (int, Assoc)  { return p.prc, p.ass }
func (p Operator) Token() token.Token   { return p.tok }

// Go lexer currently recognizes 89 tokens — not just operators — but we need
// to have enough indices for an operator token constant defined anywhere in the
// enumerated list of tokens. This enables constant-time lookup for operator
// tokens, instead of a map's relatively more-expensive hashing function.
const maxOperators = 128

type Table [maxOperators]Operator

func (m *Table) Reset() { *m = Table{} }

func (m *Table) Get(tok token.Token) (op Operator, ok bool) {
	if 0 <= int(tok) && int(tok) < maxOperators {
		op, ok = m[tok], (m[tok].bll != Unbound) || (m[tok].blr != Unbound)
	}
	return
}

func (m *Table) Add(prec int, assoc Assoc, tok ...token.Token) {
	lhs, rhs := assoc.Level(prec)
	for _, t := range tok {
		m[t] = Operator{bll: lhs, blr: rhs, prc: prec, ass: assoc, tok: t}
	}
}

type Schema struct{ prefix, postfix, infix Table }

func (s *Schema) Reset() {
	s.prefix.Reset()
	s.postfix.Reset()
	s.infix.Reset()

	s.prefix.Add(23, UnaryRight, token.LPAREN)

	s.infix.Add(22, BinaryRight, token.PERIOD)

	s.postfix.Add(21, UnaryLeft, token.INC, token.DEC)

	s.prefix.Add(20, UnaryRight, token.INC, token.DEC)
	s.prefix.Add(19, UnaryRight, token.ADD, token.SUB)
	s.prefix.Add(18, UnaryRight, token.NOT, token.TILDE)

	s.infix.Add(17, BinaryLeft, token.MUL, token.QUO, token.REM)
	s.infix.Add(16, BinaryLeft, token.ADD, token.SUB)
	s.infix.Add(15, BinaryLeft, token.SHL, token.SHR)
	s.infix.Add(14, BinaryLeft, token.LSS, token.GTR, token.LEQ, token.GEQ)
	s.infix.Add(13, BinaryLeft, token.EQL, token.NEQ)
	s.infix.Add(12, BinaryLeft, token.AND)
	s.infix.Add(11, BinaryLeft, token.AND_NOT)
	s.infix.Add(10, BinaryLeft, token.XOR)
	s.infix.Add(9, BinaryLeft, token.OR)
	s.infix.Add(8, BinaryLeft, token.LAND)
	s.infix.Add(7, BinaryLeft, token.LOR)
	s.infix.Add(6, BinaryRight, token.DEFINE, token.ASSIGN)
	s.infix.Add(5, BinaryRight, token.ADD_ASSIGN, token.SUB_ASSIGN)
	s.infix.Add(4, BinaryRight, token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN)
	s.infix.Add(3, BinaryRight, token.SHL_ASSIGN, token.SHR_ASSIGN)
	s.infix.Add(2, BinaryRight, token.AND_ASSIGN, token.AND_NOT_ASSIGN, token.XOR_ASSIGN, token.OR_ASSIGN)
	s.infix.Add(1, BinaryLeft, token.COMMA, token.SEMICOLON)

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

func (s *Schema) Prefix(t token.Token) (Operator, bool) {
	return s.prefix.Get(t)
}

func (s *Schema) Postfix(tok token.Token) (Operator, bool) {
	return s.postfix.Get(tok)
}

func (s *Schema) Infix(tok token.Token) (Operator, bool) {
	return s.infix.Get(tok)
}

var Default Schema
