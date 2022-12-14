package sym

import (
	"go/scanner"
	"go/token"
)

type scan struct {
	*scanner.Scanner
	*scanner.ErrorList
	Gate chan Symbol
	Quit chan Symbol
}

// Scan creates a new scan and a new goroutine that tokenizes the input into a
// [Symbol] channel.
//
// Each input Symbol is gated in the channel field Gate until it is received by
// another goroutine.
// Once the Symbol is received, the tokenizer will continue concurrently
// scanning for the next Symbol, and so on.
//
// Sending any Symbol to the channel field Quit forces the scanner and goroutine
// to exit immediately, even if a Symbol is currently gated.
// That Symbol would never be sent on the Gate channel and thus lost forever.
func Scan(buffer []byte) *scan {
	// Use mode = scanner.ScanComments to emit COMMENT tokens.
	const mode scanner.Mode = 0
	s := &scan{
		Scanner:   &scanner.Scanner{},
		ErrorList: &scanner.ErrorList{},
		Gate:      make(chan Symbol),
		Quit:      make(chan Symbol),
	}
	// Buffer must contain the entire input.
	// No additional input may be appended to a Streamer.
	s.Init(
		token.NewFileSet().AddFile("", -1, len(buffer)),
		buffer,
		func(pos token.Position, msg string) {
			logf("error (%s): %s", pos, msg)
			s.Add(pos, msg)
		},
		mode,
	)
	// Spawn a goroutine to repeatedly tokenize the input into the Symbol channel
	// field Gate.
	go s.run()
	return s
}

func (s *scan) run() {
	var halt bool
	for !halt {
		pos, tok, lit := s.Scan()
		u := Symbol{Token: tok, Lit: lit, Pos: pos}
		logf("Lex: SCAN: \"%s\" [%+v]", u, u)
		// Block until we have either received a Symbol on the Quit channel, or
		// something has received the Symbol being sent on the Gate channel.
		select {
		// Stop scanning input and exit this goroutine once anything is received on
		// the Quit channel.
		case q := <-s.Quit:
			logf("Lex: QUIT: \"%s\" [%+v]", q, q)
			halt = true
		// Once the Symbol is received on the Gate channel, we immediately begin
		// scanning input for the next Symbol.
		case s.Gate <- u:
			// logf("Lex: GATE (in): %s", u)
		}
	}
}
