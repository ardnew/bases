package sym

import (
	"go/scanner"
	"go/token"

	"github.com/ardnew/bases/log"
)

// Stream scans and sends the current input [Symbol] to the given channel and
// returns a Stream that will repeat this process on the next input Symbol.
//
// Stream can be considered a "state function" as described by Rob Pike in
// [Lexical Scanning in Go] (GTUG Sydney; 30 August 2011).
//
// The stateful properties of Stream are captured by a closure returned by Emit.
//
// [Lexical Scanning in Go]: https://go.dev/talks/2011/lex.slide#19
type Stream func(chan Symbol) Stream

// Emit initializes and returns a Stream ready to tokenize a given input buffer.
func Emit(buffer []byte) (stream Stream) {
	// Use mode = scanner.ScanComments to emit COMMENT tokens.
	const mode scanner.Mode = 0

	scan := &scanner.Scanner{}
	errs := &scanner.ErrorList{}
	slog := log.New(log.DefaultWriter, log.DefaultFormat)

	// The input buffer given with each call to Emit represents the entire input
	// being scanned.
	// No additional input may be appended.
	scan.Init(
		token.NewFileSet().AddFile("", -1, len(buffer)),
		buffer,
		func(pos token.Position, msg string) {
			slog.Printf("error (%s): %s", pos, msg)
			errs.Add(pos, msg)
		},
		mode,
	)

	// If the current input Symbol scanned is invalid (e.g., EOF, illegal, etc.),
	// then stop scanning and return nil.
	// Otherwise, the default next state returned from Stream is itself.
	stream = func(c chan Symbol) Stream {
		pos, tok, lit := scan.Scan()
		s := Symbol{Token: tok, Lit: lit, Pos: pos}
		c <- s
		c <- Symbol{}
		if s.IsEOF() || s.IsIllegal() {
			return nil
		}
		return stream
	}
	return
}

// Undo returns a Stream that outputs the given Symbol without scanning input.
//
// Calling the Stream result will return the receiver of Undo. This relationship
// allows for unlimited nesting, which means unlimited lookahead. For example:
//
//	...
func (e Stream) Undo(s Symbol) Stream {
	return func(c chan Symbol) Stream {
		c <- s
		return e
	}
}
