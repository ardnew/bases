package state

import (
	"github.com/ardnew/bases/io/syntax/parse/atom"
)

type Emitter interface {
	Emit(e atom.Stream)
	Undo(a atom.Atom)
}

type Parser func(e Emitter) Parser
