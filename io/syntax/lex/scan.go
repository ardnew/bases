package lex

import (
	"go/scanner"
	"go/token"
	"os"

	"github.com/ardnew/bases/io/log"
)

// Scan contains the Go lexical scanner and provides methods for receiving a
// stream of tokens from a given input string.
type Scan struct {
	dlog *log.Log
	fset *token.FileSet
	prev []Token
	seen uint
	mode scanner.Mode
	scan scanner.Scanner
}

// NewScan creates a new Scan. Users must initialize the returned Scan object by
// calling Init with the input string to tokenize. A lookbehind buffer of length
// size is allocated.
func NewScan(size int) *Scan {
	return &Scan{
		dlog: log.New(os.Stderr, log.FmtDefault),
		fset: token.NewFileSet(),
		prev: make([]Token, size),
		mode: scanner.ScanComments,
		scan: scanner.Scanner{},
	}
}

// Init initializes a Scan object by defining the input string to tokenize.
func (s *Scan) Init(src string) *Scan {
	s.scan.Init(s.fset.AddFile("", -1, len(src)), []byte(src), s.error, s.mode)
	for i := range s.prev {
		s.prev[i] = Token{}
	}
	s.seen = 0
	return s
}

// Next scans and returns the next token from the input string and then returns
// true. False is returned when no tokens remain (EOF).
func (s *Scan) Next(t *Token) bool {
	pos, tok, lit := s.scan.Scan()
	if tok == token.EOF {
		return false
	}
	*t = Token{
		Pos: s.fset.Position(pos),
		Tok: tok,
		Lit: lit,
	}
	if n := len(s.prev); n > 0 {
		m := uint(n) - 1
		if m > s.seen {
			m = s.seen
		}
		copy(s.prev[1:], s.prev[:m]) // shift the existing elements
		s.prev[0] = *t               // copy the Token, not the reference
	}
	s.seen++
	return true
}

// Prev returns the Token scanned during the n'th previous call to Next.
// The most recently scanned Token is returned when n = 0, the token before that
// is returned when n = 1, and so on.
// The size of the lookbehind buffer is defined via argument to NewScan.
func (s *Scan) Prev(n int) (t Token, ok bool) {
	if n < 0 || n >= len(s.prev) {
		return
	}
	return s.prev[n], true
}

// All calls Next repeatedly to construct and return a slice of all remaining
// tokens in the input string.
func (s *Scan) All() []Token {
	t := []Token{}
	var p Token
	for s.Next(&p) {
		t = append(t, p)
	}
	return t
}

func (s *Scan) error(pos token.Position, msg string) {
	s.dlog.Printf("error (%s): %s", pos, msg)
}
