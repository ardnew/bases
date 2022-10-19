package parse

import (
	"github.com/ardnew/bases/io/syntax/parse/lex"
)

type State interface {
	Next(lex lex.Lexer) State
	String() string
}
