package op

// Level represents the binding level between one operator and one operand.
//
// Precedence and associativity are encoded together as a single integer
// referred to throughout theese packages as a "binding level", or just "level".
// This binding level metric is constructed as follows:
//
// Every Operator has two binding levels — one level between the Operator and
// the operand on each side — left-hand side (LHS) and right-hand side (RHS).
// Unary Operators use Unbound on the side without an operand.
//
// If level a > level b, then a has higher precedence than b.
// The least binding level is defined as 0. Unbound is equivalent to 0.
// Consequently, the least level of an Operator bound to an operand is > 0.
// By convention, we choose 1 for the Operator with least precedence.
//
// To encode associativity in this metric, we only need to ensure the levels on
// either side of a given Operator are unequal, and that both of these levels
// are a.) less than both levels of every Operater with greater precedence, and
// b.) greater than both levels of every Operator with lesser precedence.
//
// A simple way to ensure this is by assigning every Operator an order of
// precedence number N > 0. The binding levels on either side of the operator
// are then computed as 2×N and 2×N±1.
//
// Since all other Operators have levels 2×M±1 != 2×N±1, this difference in
// levels on each side of the Operator only ever matters when we need a tiebreak
// among equal Operators! This realizes the associativity relationship.
//
// Parity selection depends on the particular parser being used. But for the
// Pratt parser in this program, where increasing binding level corresponds to
// increasing precedence, the associative side with higher precedence is
// assigned 2×N-1, and the other side with lower precedence 2×N.
//
// It does seems counterintuitive, admittedly, to assign the high-precedent side
// a level less than the low-precedent side's level (in a strictly increasing
// precedence metric). But this is somewhat similar to the orderings of
// production rules in a BNF grammar — the rules that appear first near the top,
// composed of other non-terminal productions, are those that are expanded last.
type Level interface {
	// Int returns the binding level as a native int type.
	Int() int
}

// Unbound represents an absent binding level, esp. for unary Operators
// associated with a single operand on either the LHS or RHS.
var Unbound Level

// lvl is a minimal implementation of Level.
type lvl int

// Int returns the binding level as a native int type.
func (i lvl) Int() int { return int(i) }
