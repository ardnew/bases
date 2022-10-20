package parse

import (
	"github.com/ardnew/bases/io/syntax/parse/lex"
)

type Parser interface {
	Parse(lex lex.Lexer) Parser
	String() string
}
