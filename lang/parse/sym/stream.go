package sym

import (
	"go/scanner"
	"go/token"

	"github.com/ardnew/bases/log"
)

// Streamer scans and sends the current input [Symbol] to the given channel and
// returns a Streamer that will repeat this process on the next input Symbol.
//
// Streamer can be considered a "state function" as described by Rob Pike in
// [Lexical Scanning in Go] (GTUG Sydney; 30 August 2011).
//
// The stateful properties of Streamer are captured by a closure returned by
// func Stream.
//
// [Lexical Scanning in Go]: https://go.dev/talks/2011/lex.slide#19
type Streamer func(Valve) Streamer

// Stream creates a Streamer ready to tokenize input from a given buffer.
func Stream(buffer []byte) (s Streamer) {
	// Use mode = scanner.ScanComments to emit COMMENT tokens.
	const mode scanner.Mode = 0

	scan := &scanner.Scanner{}
	errs := &scanner.ErrorList{}
	slog := log.New(log.DefaultWriter, log.DefaultFormat)

	// Buffer must contain the entire input.
	// No additional input may be appended to a Streamer.
	scan.Init(
		token.NewFileSet().AddFile("", -1, len(buffer)),
		buffer,
		func(pos token.Position, msg string) {
			slog.Printf("error (%s): %s", pos, msg)
			errs.Add(pos, msg)
		},
		mode,
	)

	// The Streamer s must be named so that its definition can refer to itself
	// recursively.
	s = func(v Valve) Streamer {
		// Scanner must always make progress and output the Symbol it disovered.
		pos, tok, lit := scan.Scan()
		u := Symbol{Token: tok, Lit: lit, Pos: pos}
		v.Symbol <- u
		// If the scanned input Symbol indicates end of the stream,
		// then stop scanning, close output channel, and return nil.
		if v.Close(u) {
			close(v.Symbol)
			return nil
		}
		// Otherwise, the default next state returned from Streamer is itself.
		return s
	}
	return
}

// Undo returns a Streamer that outputs the given Symbol without scanning input
// and then returns the receiver of Undo.
//
// This relationship allows for unlimited nesting, i.e., unlimited lookahead.
// For example:
//
//	...
func (s Streamer) Undo(u Symbol) Streamer {
	return func(v Valve) Streamer {
		v.Symbol <- u
		return s
	}
}

type Valve struct {
	Symbol chan Symbol
	Close  func(Symbol) bool
}

func MakeValve(f func(Symbol) bool) Valve {
	return Valve{Symbol: make(chan Symbol), Close: f}
}
