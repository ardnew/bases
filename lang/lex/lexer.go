package lex

import (
	"github.com/ardnew/bases/lang/lex/sym"
)

// Lexer defines the interface for processing a stream of tokens read from an
// input string or byte source.
//
// Conceptually, there are two driving processes:
//
//  1. Lexing tokens – generates indivisble grammar component (token) stream
//  2. Parsing tokens – verifies/evaluates token stream's syntax and semantics
//
// Since each process may introduce overhead or latency independent of the other
// process, they are good candidates for a concurrent design.
type Lexer interface {
	Take() sym.Symbol
	TakeNext([]sym.Symbol) error
	Untake(...sym.Symbol)
	Look() sym.Symbol
	Check(sym.Symbol) bool
	Err() error
	Error() string
}
