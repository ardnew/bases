package op

import "go/token"

type Table struct {
	prefix  operators
	postfix operators
	infix   operators
}

var Default Table

func (t *Table) Prefix(tok token.Token) (Operator, bool) {
	return t.prefix.get(tok)
}

func (t *Table) Postfix(tok token.Token) (Operator, bool) {
	return t.postfix.get(tok)
}

func (t *Table) Infix(tok token.Token) (Operator, bool) {
	return t.infix.get(tok)
}

func (s *Table) Reset() {
	*s = Table{}

	s.prefix.add(23, UnaryRight, token.LPAREN)
	s.infix.add(22, BinaryRight, token.PERIOD)
	s.postfix.add(21, UnaryLeft, token.INC, token.DEC)
	s.prefix.add(20, UnaryRight, token.INC, token.DEC)
	s.prefix.add(19, UnaryRight, token.ADD, token.SUB)
	s.prefix.add(18, UnaryRight, token.NOT, token.TILDE)
	s.infix.add(17, BinaryLeft, token.MUL, token.QUO, token.REM)
	s.infix.add(16, BinaryLeft, token.ADD, token.SUB)
	s.infix.add(15, BinaryLeft, token.SHL, token.SHR)
	s.infix.add(14, BinaryLeft, token.LSS, token.GTR, token.LEQ, token.GEQ)
	s.infix.add(13, BinaryLeft, token.EQL, token.NEQ)
	s.infix.add(12, BinaryLeft, token.AND)
	s.infix.add(11, BinaryLeft, token.AND_NOT)
	s.infix.add(10, BinaryLeft, token.XOR)
	s.infix.add(9, BinaryLeft, token.OR)
	s.infix.add(8, BinaryLeft, token.LAND)
	s.infix.add(7, BinaryLeft, token.LOR)
	s.infix.add(6, BinaryRight, token.DEFINE, token.ASSIGN)
	s.infix.add(5, BinaryRight, token.ADD_ASSIGN, token.SUB_ASSIGN)
	s.infix.add(4, BinaryRight, token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN)
	s.infix.add(3, BinaryRight, token.SHL_ASSIGN, token.SHR_ASSIGN)
	s.infix.add(2, BinaryRight, token.AND_ASSIGN, token.AND_NOT_ASSIGN, token.XOR_ASSIGN, token.OR_ASSIGN)
	s.infix.add(1, BinaryLeft, token.COMMA, token.SEMICOLON)

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

func init() { Default.Reset() }
