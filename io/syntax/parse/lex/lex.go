package lex

import (
	"go/token"

	"github.com/ardnew/bases/io/syntax/parse/lex/atom"
)

type Lexer interface {
	Emit() (a atom.Atom)
	Peek() (a atom.Atom)
	Undo(a atom.Atom)
	Fail(a atom.Atom, expect ...token.Token)
}
