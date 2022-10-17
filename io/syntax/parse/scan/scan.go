package scan

import (
	"go/scanner"
	"go/token"
	"os"

	"github.com/ardnew/bases/io/log"
	"github.com/ardnew/bases/io/syntax/parse/atom"
	"github.com/ardnew/bases/io/syntax/parse/state"
)

// Scan contains the Go lexical scanner and provides methods for progressively
// emitting tokens from the input stream.
type Scan struct {
	log   *log.Log
	state state.Parser
	undo  atom.Stream
	fset  *token.FileSet
	errs  scanner.ErrorList
	scanner.Scanner
}

// NewScan creates a new scan. Users must initialize the returned scan object by
// calling init with the input string to tokenize.
func NewScan() *Scan {
	return &Scan{
		log:     log.New(os.Stderr, log.FmtDefault),
		undo:    make(atom.Stream),
		fset:    token.NewFileSet(),
		Scanner: scanner.Scanner{},
	}
}

// init initializes a scan object by defining the input string to tokenize.
func (s *Scan) Init(src string) *Scan {
	const mode = 0 // Use scanner.ScanComments to emit COMMENT tokens.
	s.errs.Reset()
	file := s.fset.AddFile("", -1, len(src))
	s.Scanner.Init(file, []byte(src), s.addError, mode)
	return s
}

func (s *Scan) addError(pos token.Position, msg string) {
	s.log.Printf("error (%s): %s", pos, msg)
	s.errs.Add(pos, msg)
}

func (s *Scan) Error() (err string) {
	if s.errs.Len() > 0 {
		s.errs.RemoveMultiples()
		err = s.errs.Error()
	}
	return
}

func (s *Scan) Run() error {
	for {
		s.state = s.state(s)
	}
}

// Emit scans the next token from the input stream and sends it to the input
// channel read if not EOF, otherwise the channel is closed immediately.
func (s *Scan) Emit(e atom.Stream) {
	var next atom.Atom
	select {
	case next = <-s.undo:
	default:
		pos, tok, lit := s.Scan()
		next = atom.Atom{Lit: lit, Token: tok, Pos: pos}
	}
	if next.Token == token.EOF {
		close(e)
	} else {
		e <- next
	}
}

func (s *Scan) Undo(a atom.Atom) { s.undo <- a }
