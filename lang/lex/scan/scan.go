package scan

import (
	"go/scanner"
	"go/token"
	"strings"

	"github.com/ardnew/bases/lang/lex/sym"
	"github.com/ardnew/bases/list"
	"github.com/ardnew/bases/log"
)

// Scan contains the Go lexical scanner and provides methods for progressively
// emitting tokens from the input stream.
type Scan struct {
	log  *log.Log
	fset *token.FileSet
	stop chan interface{}
	next sym.Stream
	undo list.SyncStack
	errs scanner.ErrorList
}

// New creates a new Scan.
func New() *Scan { return (&Scan{}).Reset() }

// Reset initializes all fields to their default state.
//
// All line information previously recorded in the internal file set is cleared
// along with any errors generated when scanning it.
//
// The output channel and lookahead buffer are both flushed.
func (s *Scan) Reset() *Scan {
	if s == nil {
		return New()
	}
	s.log = log.New(log.DefaultWriter, log.DefaultFormat)
	s.fset = token.NewFileSet()
	s.stop = make(chan interface{}, 1)
	if s.next != nil {
		close(s.next)
	}
	s.next = make(sym.Stream)
	s.undo = list.SyncStack{}
	s.errs.Reset()
	return s
}

// AddBuffer appends a buffer to the scanner input amd begins processing it
// immediately in a new goroutine.
func (s *Scan) AddBuffer(b []byte) error {
	file := s.fset.AddFile("", -1, len(b))
	go s.emit(file, b, token.ILLEGAL, token.EOF)
	return s
}

// emit continuously scans tokens from the input source and sends them to the
// output channel.
// It stops scanning after any given [token.Token] in eos is scanned and sent,
// which should probably always include [token.ILLEGAL] and [token.EOF].
//
// emit is intended to be run in its own goroutine. This allows the lexer to
// proceed over the input while a parser concurrently processes the tokens.
func (s *Scan) emit(f *token.File, b []byte, eos ...token.Token) {
	const mode = 0 // Use scanner.ScanComments to emit COMMENT tokens.
	sc := &scanner.Scanner{}
	sc.Init(f, b, s.addError, mode)
	for {
		pos, tok, lit := sc.Scan()
		s.next <- sym.Symbol{Lit: lit, Token: tok, Pos: pos}
		for _, e := range eos {
			if tok == e {
				return
			}
		}
	}
}

func (s *Scan) Take() sym.Symbol {
	if a, ok := s.undo.Pop().(sym.Symbol); ok {
		return a
	}
	return <-s.next
}

func (s *Scan) Untake(a ...sym.Symbol) {
	for _, t := range a {
		s.undo.Push(t)
	}
}

func (s *Scan) Look() sym.Symbol {
	a := s.Take()
	s.Untake(a)
	return a
}

func (s *Scan) Error() (err string) {
	if s.errs.Len() > 0 {
		s.errs.RemoveMultiples()
		err = s.errs.Error()
	}
	return
}

func (s *Scan) Err() (err error) {
	if s.errs.Len() > 0 {
		return s
	}
	return nil
}

func (s *Scan) Fail(a sym.Symbol, expect ...token.Token) {
	var b strings.Builder
	b.WriteString("invalid token: ")
	b.WriteString(a.String())
	if n := len(expect); n > 0 {
		b.WriteString(" (expected: ")
		switch n {
		case 1:
			b.WriteString(sym.QQ(expect[0]))
		case 2:
			b.WriteString(sym.QQ(expect[0]))
			b.WriteString(" or ")
			b.WriteString(sym.QQ(expect[1]))
		default:
			for i, e := range expect {
				if i > 1 {
					b.WriteString(", ")
					if i == n {
						b.WriteString("or ")
					}
				}
				b.WriteString(sym.QQ(e))
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
