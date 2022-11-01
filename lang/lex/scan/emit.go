package scan

import (
	"go/scanner"
	"go/token"

	"github.com/ardnew/bases/lang/lex/sym"
	"github.com/ardnew/bases/log"
)

// Emit scans and sends the next [sym.Symbol] from its input to the given
// [sym.Stream], and it returns an Emit that will repeat this process on the
// Symbol following that.
//
// The stateful properties of Emit are defined and captured by a closure over
// local variables in function Make.
type Emit func(sym.Stream) Emit

// Make initializes and returns an Emit state function ready to tokenize the
// given input buffer.
func Make(buff []byte) (emit Emit) {
	// Use mode = scanner.ScanComments to emit COMMENT tokens.
	const mode scanner.Mode = 0

	// Create a new Scanner instance for doing all the heavy lifting such as
	// lexing and tokenizing the input.
	scan := &scanner.Scanner{}
	errs := &scanner.ErrorList{}
	slog := log.New(log.DefaultWriter, log.DefaultFormat)

	// Init prepares the Scanner instance to begin tokening the input.
	//
	// The input buffer given with each call to Make represents the entire input
	// being scanned.
	// No additional input may be appended.
	// Therefore, we do not need to maintain or share the [token.FileSet],
	// and it will always consist of exactly one [token.File]: the input buffer.
	scan.Init(
		token.NewFileSet().AddFile("", -1, len(buff)),
		buff,
		func(pos token.Position, msg string) {
			slog.Printf("error (%s): %s", pos, msg)
			errs.Add(pos, msg)
		},
		mode,
	)

	// Define the state function such that each time it is called, a single token
	// will be scanned and sent to the given output channel.
	//
	// The next state function called should be the return value of its previous
	// invocation, or the function returned by a method with it as the receiver.
	emit = func(ss sym.Stream) Emit {
		pos, tok, lit := scan.Scan()
		sy := sym.Symbol{Token: tok, Lit: lit, Pos: pos}
		ss <- sy
		if sy.IsEOF() || sy.IsIllegal() {
			return nil
		}
		return emit
	}
	return
}

// Undo returns a state function that will send the given Symbol to the output
// channel without scanning a new token from the input stream.
//
// The state function returned when calling the function returned from Undo is
// the receiver of Undo. This relationship allows for unlimited nesting, which
// also means unlimited forward lookahead. For example:
//
//	sy := <-ch        // var ch sym.Stream
//	e1 := e0.Undo(sy) // e1 wraps e0 with a simple channel send
//	e2 := e1.Undo(sy) // e2 wraps e1, and e1 wraps e0
//	e3 := e2(ch)      // calling e2 returns the e1 wrapper over e0
//	e4 := e3(ch)      // calling e3 returns the original e0 unmodified
func (e Emit) Undo(s sym.Symbol) Emit {
	return func(ss sym.Stream) Emit {
		ss <- s
		return e
	}
}
