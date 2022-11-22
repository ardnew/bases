package sym

// Streamer returns the current input [Symbol] and a Streamer that will return
// the next input Symbol.
//
// Streamer can be considered a "state function" as described by Rob Pike in
// [Lexical Scanning in Go] (GTUG Sydney; 30 August 2011).e
//
// The stateful properties of Streamer are captured by a closure returned by
// func Stream.
//
// [Lexical Scanning in Go]: https://go.dev/talks/2011/lex.slide#19
type Streamer func() (Symbol, Streamer)

// Stopper returns true if the given [Symbol] represents an end of file
// (or end of stream). It returns false otherwise.
type Stopper func(Symbol) bool

// Stream creates a Streamer ready to tokenize input from a given buffer.
//
// If any given [Stopper] function returns true after scanning an input Symbol,
// then that Symbol is returned with a nil Streamer.
// That Symbol is also sent to the Quit channel to halt its concurrent Scan.
func Stream(buffer []byte, p ...Stopper) (s Streamer) {
	scan := Scan(buffer)
	// The Streamer s must be named so that its definition can refer to itself
	// recursively.
	s = func() (Symbol, Streamer) {
		u := <-scan.Gate
		logf("Stream: GATE: %s", u)
		// If the scanned Symbol indicates the end of input,
		// then halt the scanner, close output channel, and return nil.
		for _, e := range p {
			if e(u) {
				logf("Stream: STOP: \"%s\" [%+v]", u, u)
				scan.Quit <- u
				return u, nil
			}
		}
		// Reuse the closure; Scanner is positioned immediately after the unread
		// Symbol in the output channel.
		return u, s
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
	if s == nil {
		return func() (Symbol, Streamer) { return EOF(), nil }
	}
	return func() (Symbol, Streamer) { return u, s }
}

// Next returns the current Symbol in the input stream and begins scanning for
// the next.
func (s *Streamer) Next() (a Symbol) {
	if *s == nil {
		return EOF()
	}
	a, *s = (*s)()
	return
}

// Peek returns the current Symbol in the input stream.
func (s *Streamer) Peek() (a Symbol) {
	a = s.Next()
	*s = s.Undo(a)
	return
}
