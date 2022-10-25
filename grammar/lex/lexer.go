package lex

import (
	"go/token"

	"github.com/ardnew/bases/grammar/lex/sym"
)

type Lexer interface {
	Take() sym.Symbol
	TakeNext(a []sym.Symbol) error
	Untake(a ...sym.Symbol)
	Look() sym.Symbol
	Check(a sym.Symbol) (ok bool)
	Fail(a sym.Symbol, expect ...token.Token)
}
