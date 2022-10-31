package oper

// Assoc represents an operator's associativity.
//
// Associativity is an operator's tendency to bind to an operand on the left-
// hand side (LHS) or right-hand side (RHS) when evaluating expressions
// containing operators of equal precedence.
type Assoc int

const (
	Nonassociative Assoc = iota
	UnaryLeft
	UnaryRight
	BinaryLeft
	BinaryRight
	assocCount
)

// Level returns the LHS and RHS binding levels of an s-associative operator
// with given precedence.
func (s Assoc) Level(precedence int) (lhs, rhs Level) {
	if precedence > 0 {
		n := 2 * precedence
		switch s {
		case UnaryLeft:
			return lvl(n - 1), Unbound
		case UnaryRight:
			return Unbound, lvl(n - 1)
		case BinaryLeft:
			return lvl(n - 1), lvl(n)
		case BinaryRight:
			return lvl(n), lvl(n - 1)
		}
	}
	return Unbound, Unbound
}
