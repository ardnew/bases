package op

import (
	"go/token"
	"sort"
)

var Prefix, Postfix, Infix Table

func init() {
	Prefix.Reset()
	Postfix.Reset()
	Infix.Reset()

	Prefix.Add(23, UnaryRight, token.LPAREN)

	Infix.Add(22, BinaryRight, token.PERIOD)

	Postfix.Add(21, UnaryLeft, token.INC, token.DEC)

	Prefix.Add(20, UnaryRight, token.INC, token.DEC)
	Prefix.Add(19, UnaryRight, token.ADD, token.SUB)
	Prefix.Add(18, UnaryRight, token.NOT, token.TILDE)

	Infix.Add(17, BinaryLeft, token.MUL, token.QUO, token.REM)
	Infix.Add(16, BinaryLeft, token.ADD, token.SUB)
	Infix.Add(15, BinaryLeft, token.SHL, token.SHR)
	Infix.Add(14, BinaryLeft, token.LSS, token.GTR, token.LEQ, token.GEQ)
	Infix.Add(13, BinaryLeft, token.EQL, token.NEQ)
	Infix.Add(12, BinaryLeft, token.AND)
	Infix.Add(11, BinaryLeft, token.AND_NOT)
	Infix.Add(10, BinaryLeft, token.XOR)
	Infix.Add(9, BinaryLeft, token.OR)
	Infix.Add(8, BinaryLeft, token.LAND)
	Infix.Add(7, BinaryLeft, token.LOR)
	Infix.Add(6, BinaryRight, token.DEFINE, token.ASSIGN)
	Infix.Add(5, BinaryRight, token.ADD_ASSIGN, token.SUB_ASSIGN)
	Infix.Add(4, BinaryRight, token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN)
	Infix.Add(3, BinaryRight, token.SHL_ASSIGN, token.SHR_ASSIGN)
	Infix.Add(2, BinaryRight, token.AND_ASSIGN, token.AND_NOT_ASSIGN, token.XOR_ASSIGN, token.OR_ASSIGN)
	Infix.Add(1, BinaryLeft, token.COMMA, token.SEMICOLON)

	Prefix.Sort()
	Postfix.Sort()
	Infix.Sort()

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

// Operator provides an abstraction for any type of operator in prefix, postfix,
// infix, or other complex expressions.
type Operator struct {
	L, R Level
	token.Token
}

func (p Operator) String() string { return p.Token.String() }

func (p Operator) precedence() (n int) {
	defer func() { n /= 2 }()
	switch {
	case p.L == Unbound && p.R == Unbound:
		return 0
	case p.L == Unbound:
		return p.R.Int() + 1
	case p.R == Unbound:
		return p.L.Int() + 1
	default:
		l, r := p.L.Int(), p.R.Int()
		if l < r {
			return r
		}
		return l
	}
}

type Table struct {
	Lut map[token.Token]Operator // Operator lookup table, keyed by Token.
	Ord []Operator               // All elements in Op sorted by precedence.
}

func (m *Table) Reset() {
	if m != nil {
		*m = Table{
			Lut: map[token.Token]Operator{},
			Ord: []Operator{},
		}
	}
}

// Add adds an arbitrary number of given [token.Token] keys mapped to Operator
// values constructed with the given precedence, associativity, and
// [state.Parser] state function.
//
// You must call Sort after all elements have been added via method Add.
// Otherwise, the elements in Ord may not be ordered by precedence.
func (m *Table) Add(precedence int, assoc Assoc, tok ...token.Token) {
	lhs, rhs := assoc.Level(precedence)
	for _, t := range tok {
		m.Lut[t] = Operator{Token: t, L: lhs, R: rhs}
		m.Ord = append(m.Ord, m.Lut[t])
	}
}

// Sort sorts the elements of Ord based on precedence, while preserving the
// order of any elements with equal precedence.
// Sort does not read or modify Lut.
//
// Should be called only once after all elements have been added to the table
// via method Add.
func (m *Table) Sort() { sort.Stable(m) }

// Len is the number of elements in the collection.
func (m *Table) Len() int { return len(m.Ord) }

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//   - if both Less(i, j) and Less(j, k) are true,
//     then Less(i, k) must be true as well.
//   - if both Less(i, j) and Less(j, k) are false,
//     then Less(i, k) must be false as well.
//
// Note that floating-point comparison (operator < on float32/float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (m *Table) Less(i, j int) bool {
	return m.Ord[i].precedence() > m.Ord[j].precedence()
}

// Swap swaps the elements with indexes i and j.
func (m *Table) Swap(i, j int) {
	o := m.Ord[i]
	m.Ord[i] = m.Ord[j]
	m.Ord[j] = o
}
