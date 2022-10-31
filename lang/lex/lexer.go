package lex

import (
	"github.com/ardnew/bases/lang/lex/sym"
)

// Lexer defines the interface for processing a stream of tokens read from an
// input string or byte source.
//
// Conceptually, there are two driving processes:
//
//  1. Lexer – generates stream of indivisble grammar components (tokens)
//  2. Parser – verifies/evaluates syntax and semantics of token stream
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
