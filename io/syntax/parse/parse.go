package parse

import (
	"github.com/ardnew/bases/io/syntax/parse/lex"
)

type State func(lex lex.Lexer) State
