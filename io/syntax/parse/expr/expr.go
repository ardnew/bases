package expr

import (
	"strings"

	"github.com/ardnew/bases/io/syntax/parse"
	"github.com/ardnew/bases/io/syntax/parse/expr/op"
	"github.com/ardnew/bases/io/syntax/parse/lex"
	"github.com/ardnew/bases/io/syntax/parse/lex/atom"
)

type Expr interface {
	String() string
	Parse(lex lex.Lexer) parse.State
}

// var literal = []token.Token{
// 	token.IDENT, token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING,
// }

func Parse(lex lex.Lexer) (state parse.State) {
	return (&Prefix{}).Parse
}

// func parsePrim(lex lex.Lexer) (e Prim) {
// }

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

// Parse parses tokens from the given Lexer to construct this expression and
// returns a parser for the next valid tokens.
// Nil is returned if the Lexer does not emit valid tokens for this expression.
func (x *Prim) Parse(lex lex.Lexer) parse.State {
	a := lex.Emit()
	if a.IsLiteral() {
		x.Atom = a
		return Parse
	}
	go lex.Undo(a)
	return nil

// Ident represents a variable identifier following the definition from the Go
// language spec (and implemented by the Go lexical scanner).
type Ident struct{ atom.Atom }

func (x *Ident) String() string {
	return x.Atom.String()
}

// Parse parses tokens from the given Lexer to construct this expression and
// returns a parser for the next valid tokens.
// Nil is returned if the Lexer does not emit valid tokens for this expression.
func (x *Ident) Parse(lex lex.Lexer) parse.State {
	return Parse
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
// The error above is because the assignment operator "=" is right-associative,
// which results in attempting to evaluate "11 = 3", and even this weirdo
// program doesn't allow assignment to numeric literals — an Ident expression is
// the only valid expression on the LHS of an assignment.
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

// Parse parses tokens from the given Lexer to construct this expression and
// returns a parser for the next valid tokens.
// Nil is returned if the Lexer does not emit valid tokens for this expression.
func (x *Assign) Parse(lex lex.Lexer) parse.State {
	return Parse
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

// Parse parses tokens from the given Lexer to construct this expression and
// returns a parser for the next valid tokens.
// Nil is returned if the Lexer does not emit valid tokens for this expression.
func (x *Infix) Parse(lex lex.Lexer) parse.State {
	return Parse
}

type Prefix struct {
	Expr
	op.Operator
}

func (x *Prefix) String() string {
	var b strings.Builder
	b.WriteString(" (")
	b.WriteString(x.Operator.String())
	b.WriteString(x.Expr.String())
	b.WriteString(") ")
	return b.String()
}

// Parse parses tokens from the given Lexer to construct this expression and
// returns a parser for the next valid tokens.
// Nil is returned if the Lexer does not emit valid tokens for this expression.
func (x *Prefix) Parse(lex lex.Lexer) parse.State {
	for _, o := range op.Prefix.Ord {

	}
	return Parse
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

// Parse parses tokens from the given Lexer to construct this expression and
// returns a parser for the next valid tokens.
// Nil is returned if the Lexer does not emit valid tokens for this expression.
func (x *Postfix) Parse(lex lex.Lexer) parse.State {
	return Parse
}
