package sym

import (
	"go/scanner"
	"go/token"
	"io"

	"github.com/ardnew/bases/log"
)

// Streamer scans and sends the current input [Symbol] to the given channel and
// returns a Streamer that will repeat this processing on the next input Symbol.
//
// If any given [Stopper] function returns true after scanning an input Symbol,
// then it becomes the final Symbol added before closing the output channel.
//
// Streamer can be considered a "state function" as described by Rob Pike in
// [Lexical Scanning in Go] (GTUG Sydney; 30 August 2011).e
//
// The stateful properties of Streamer are captured by a closure returned by
// func Stream.
//
// [Lexical Scanning in Go]: https://go.dev/talks/2011/lex.slide#19
type Streamer func() (Streamer, Symbol)

// Stopper returns true if the given [Symbol] represents an end of file
// (or end of stream). It returns false otherwise.
type Stopper func(Symbol) bool

// Stream creates a Streamer ready to tokenize input from a given buffer.
func Stream(buffer []byte, p ...Stopper) (s Streamer) {
	// Use mode = scanner.ScanComments to emit COMMENT tokens.
	const mode scanner.Mode = 0

	scan := &scanner.Scanner{}
	errs := &scanner.ErrorList{}
	slog := log.New(io.Discard, "")
	if w, f, e := log.LookupEnv("STREAM"); e == nil {
		slog = log.New(w, f)
		slog.SetCallerOffset(1)
	}

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

	gate := make(chan Symbol)
	quit := make(chan Symbol)
	go func() {
		var halt bool
		for !halt {
			pos, tok, lit := scan.Scan()
			u := Symbol{Token: tok, Lit: lit, Pos: pos}
			slog.Printf("Scan: SCAN: \"%s\" [%+v]", u, u)
			select {
			case q := <-quit:
				slog.Printf("Scan: QUIT: \"%s\" [%+v]", q, q)
				halt = true
			case gate <- u:
				slog.Printf("Scan: SEND: %s", u)
			}
		}
	}()

	// The Streamer s must be named so that its definition can refer to itself
	// recursively.
	s = func() (Streamer, Symbol) {
		u := <-gate
		slog.Printf("Stream: RECV: %s", u)
		// If the scanned Symbol indicates the end of input,
		// then halt the scanner, close output channel, and return nil.
		for _, e := range p {
			if e(u) {
				slog.Printf("Stream: STOP: \"%s\" [%+v]", u, u)
				quit <- u
				return nil, u
			}
		}
		// Reuse the closure; Scanner is positioned immediately after the unread
		// Symbol in the output channel.
		return s, u
	}
	return
}

// Undo returns a Streamer that outputs the given Symbol without scanning input
// and then returns the receiver of Undo.
//
// This allows for unlimited nesting, i.e., unlimited lookahead. For example:
//
//	...
func (s Streamer) Undo(u Symbol) Streamer {
	return func() (Streamer, Symbol) { return s, u }
}

func (s *Streamer) Next() (a Symbol) {
	*s, a = (*s)()
	return
}

// Peek returns the next Symbol in the input stream.
// Unlike receiving from a channel, it does not remove it from the stream.
//
// Peek is not technically a state function, as it does not return a Streamer.
// However, it does modify its receiver by wrapping Undo over it.
func (s *Streamer) Peek() (a Symbol) {
	a = s.Next()
	*s = s.Undo(a)
	return
}
