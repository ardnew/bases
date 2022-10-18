package scan

import (
	"go/scanner"
	"go/token"
	"os"
	"strings"

	"github.com/ardnew/bases/io/log"
	"github.com/ardnew/bases/io/syntax/parse/lex/atom"
)

// Scan contains the Go lexical scanner and provides methods for progressively
// emitting tokens from the input stream.
type Scan struct {
	log  *log.Log
	undo atom.Stream
	fset *token.FileSet
	errs scanner.ErrorList
	scanner.Scanner
}

// NewScan creates a new Scan. Users must initialize the returned Scan object by
// calling Init with the source bytes to tokenize.
func NewScan() *Scan {
	return &Scan{
		log:     log.New(os.Stderr, log.FmtDefault),
		undo:    make(atom.Stream),
		fset:    token.NewFileSet(),
		Scanner: scanner.Scanner{},
	}
}

// Init initializes a Scan object by defining the source bytes to tokenize.
func (s *Scan) Init(src []byte) *Scan {
	const mode = 0 // Use scanner.ScanComments to emit COMMENT tokens.
	s.errs.Reset()
	file := s.fset.AddFile("", -1, len(src))
	s.Scanner.Init(file, src, s.addError, mode)
	return s
}

func (s *Scan) Error() (err string) {
	if s.errs.Len() > 0 {
		s.errs.RemoveMultiples()
		err = s.errs.Error()
	}
	return
}

// Emit scans and returns the next token from the input stream.
func (s *Scan) Emit() (a atom.Atom) {
	select {
	case a = <-s.undo:
	default:
		pos, tok, lit := s.Scan()
		a = atom.Atom{Lit: lit, Token: tok, Pos: pos}
	}
	return
}

func (s *Scan) Peek() (a atom.Atom) {
	a = s.Emit()
	go s.Undo(a)
	return
}

func (s *Scan) Undo(a atom.Atom) { s.undo <- a }

func (s *Scan) Fail(a atom.Atom, expect ...token.Token) {
	var b strings.Builder
	b.WriteString("invalid token: ")
	b.WriteString(a.String())
	if n := len(expect); n > 0 {
		b.WriteString(" (expected: ")
		switch n {
		case 1:
			b.WriteString(atom.Atom{Token: expect[0]}.QQ())
		case 2:
			b.WriteString(atom.Atom{Token: expect[0]}.QQ())
			b.WriteString(" or ")
			b.WriteString(atom.Atom{Token: expect[1]}.QQ())
		default:
			for i, e := range expect {
				if i > 1 {
					b.WriteString(", ")
					if i == n {
						b.WriteString("or ")
					}
				}
				b.WriteString(atom.Atom{Token: e}.QQ())
			}
		}
		b.WriteRune(')')
	}
	s.addError(s.fset.Position(a.Pos), b.String())
}

func (s *Scan) addError(pos token.Position, msg string) {
	s.log.Printf("error (%s): %s", pos, msg)
	s.errs.Add(pos, msg)
}
