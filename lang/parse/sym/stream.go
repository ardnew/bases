package sym

import (
	"go/scanner"
	"go/token"

	"github.com/ardnew/bases/log"
)

// Streamer scans and sends the current input [Symbol] to the given channel and
// returns a Streamer that will repeat this processing on the next input Symbol.
//
// If any given [Stopper] function returns true after scanning an input Symbol,
// then it becomes the final Symbol added before closing the output channel.
//
// Streamer can be considered a "state function" as described by Rob Pike in
// [Lexical Scanning in Go] (GTUG Sydney; 30 August 2011).
//
// The stateful properties of Streamer are captured by a closure returned by
// func Stream.
//
// [Lexical Scanning in Go]: https://go.dev/talks/2011/lex.slide#19
type Streamer func(chan Symbol, ...Stopper) Streamer

// Stopper returns true if the given [Symbol] represents an end of file
// (or end of stream). It returns false otherwise.
type Stopper func(Symbol) bool

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
	s = func(c chan Symbol, p ...Stopper) Streamer {
		// Scanner must always make progress and output the Symbol it disovered.
		pos, tok, lit := scan.Scan()
		u := Symbol{Token: tok, Lit: lit, Pos: pos}
		c <- u
		// If the scanned Symbol indicates the end of input,
		// then return nil to encourage the caller to stop scanning.
		for _, e := range p {
			if e(u) {
				// Close to let other goroutines 'range' over the output channel.
				close(c)
				return nil
			}
		}
		// Reuse the closure; Scanner is positioned immediately after the unread
		// Symbol in the output channel.
		return s
	}
	return
}

// Go creates a new goroutine that scans input Symbols and sends them to the
// given channel until no input remains or a given [Stopper] returned true.
func (s *Streamer) Go(c chan Symbol, p ...Stopper) {
	go func(s *Streamer, c chan Symbol, p ...Stopper) {
		for s != nil {
			*s = (*s)(c, p...)
		}
	}(s, c, p...)
}

// Undo returns a Streamer that outputs the given Symbol without scanning input
// and then returns the receiver of Undo.
//
// This allows for unlimited nesting, i.e., unlimited lookahead. For example:
//
//	...
func (s Streamer) Undo(u Symbol) Streamer {
	return func(c chan Symbol, _ ...Stopper) Streamer {
		c <- u
		return s
	}
}

// Peek returns the next Symbol in the input stream.
// Unlike receiving from a channel, it does not remove it from the stream.
//
// Peek is not technically a state function, as it does not return a Streamer.
// However, it does modify its receiver by wrapping Undo over it.
func (s *Streamer) Peek(c chan Symbol) (a Symbol) {
	if u, ok := <-c; ok {
		a = u
		*s = s.Undo(u)
	}
	return
}
