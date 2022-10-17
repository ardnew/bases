package expr

import (
	"strings"

	"github.com/ardnew/bases/io/syntax/parse/atom"
	"github.com/ardnew/bases/io/syntax/parse/op"
	"github.com/ardnew/bases/io/syntax/parse/state"
)

type Expr interface {
	String() string
}

// Prim represents a literal in any one of the five basic types: integer,
// floating-point, imaginary, rune, or string.
//
// Prim is an abbreviation for "primitive", a catch-all term used specifically
// here because it is absent from the Go language spec, and it is thus free for
// us to define (i.e., butcher) however and wherever it makes most sense.
//
// Atom.string will hold the token literal spelling, and Atom.Token will
// distinguish the type of the token.
type Prim struct{ atom.Atom }

func (x *Prim) String() string {
	return x.Atom.String()
}

// Ident represents a variable identifier following the definition from the Go
// language spec (and implemented by the Go lexical scanner).
type Ident struct{ atom.Atom }

func (x *Ident) String() string {
	return x.Atom.String()
}

// Assign represents an assignment expression. That is, a variable identifier is
// assigned the value of an expression, and the entire assignment expression
// itself evaluates to that same value.
//
// In Go, assignments are always statements, not expressions. I find it more
// convenient to use assignments as expressions in the context of a run-eval-
// print-loop (REPL), such as this program.
//
// For example, "a = 1 + (b = 3)" is an illegal expression in Go, because the
// substring "b = 3" is a statement, not an operand in an expression. Other
// languages, including this program, permit assignment expressions by using the
// value assigned ("3", in this example) as the value of the assignment "b = 3".
// Thus, the above expression would reduce as follows:
//
// With parentheses:
//  1. a = 1 + (b = 3)
//  2. a = 1 + 3
//  3. a = 4
//
// Without parentheses (assume b = 10):
//  1. a = 1 + b = 3
//  2. a = 11 = 3      <-- error!
//
// The error on line 2 of the example above is because the assignment operator
// "=" is right-associative, which results in attempting to evaluate "11 = 3",
// and even this weirdo program doesn't allow assignment to numeric literals.
type Assign struct {
	Expr
	Ident // Identifier being assigned with Expr
}

func (x *Assign) String() string {
	var b strings.Builder
	b.WriteString(" (")
	b.WriteString(x.Ident.String())
	b.WriteRune('=')
	b.WriteString(x.Expr.String())
	b.WriteString(") ")
	return b.String()
}

type Infix struct {
	L, R Expr
	op.Operator
}

func (x *Infix) String() string {
	var b strings.Builder
	b.WriteString(" (")
	b.WriteString(x.L.String())
	b.WriteString(x.Operator.String())
	b.WriteString(x.R.String())
	b.WriteString(") ")
	return b.String()
}

type Prefix struct {
	Expr
	op.Operator
}

func (x *Prefix) Parse(e state.Emitter) state.Parser {
	m := make(atom.Stream)
	go e.Emit(m)
	a, ok := <-m
	if !ok {
		// EOF
	}
	o, ok := op.Prefix[a.Token]
	if !ok {
		// Not a prefix operator
	}
}

func (x *Prefix) String() string {
	var b strings.Builder
	b.WriteString(" (")
	b.WriteString(x.Operator.String())
	b.WriteString(x.Expr.String())
	b.WriteString(") ")
	return b.String()
}

type Postfix struct {
	Expr
	op.Operator
}

func (x *Postfix) String() string {
	var b strings.Builder
	b.WriteString(" (")
	b.WriteString(x.Expr.String())
	b.WriteString(x.Operator.String())
	b.WriteString(") ")
	return b.String()
}
