package eval

import (
	"io"

	"github.com/ardnew/bases/lang/parse"
)

type Value struct {
	*parse.Expr
}

func New() *Value {
	return &Value{Expr: parse.New()}
}

func (v *Value) Read(p []byte) (n int, err error) {
	if err = v.Err(); err == nil {
		return copy(p, v.String()), nil
	}
	return
}

func (v *Value) WriteTo(w io.Writer) (n int64, err error) {
	if err = v.Err(); err == nil {
		m, e := w.Write([]byte(v.String()))
		return int64(m), e
	}
	return
}

// Write implements the [io.Writer] interface by parsing the given buffer into
// the receiver's underlying [parse.Expr].
//
// It returns the number of bytes consumed into Expr and a nil error on success.
// Otherwise, it returns 0 and the first error encountered.
//
// It always returns either err == nil or n == 0; but it never returns both.
//
// Callers should be sure to retain unused trailing bytes on success.
// In particular, when 0 < n < len(p), err will be nil and the slice p[n:] may
// contain valid bytes for subsequent Exprs.
func (v *Value) Write(p []byte) (n int, err error) {
	m, e := v.ParseBuffer(p)
	return int(m), e
}

func (v *Value) ReadFrom(r io.Reader) (n int64, err error) {
	return v.Parse(r)
}

func (v *Value) String() string {
	if err := v.Err(); err != nil {
		return err.Error()
	}
	return v.Expr.String()
}

// func EvalString(s string) error {
// e := eval{Scan: scan.New().Init([]byte(s))}
// go e.Until(token.EOF, token.ILLEGAL)
// return e.Parse()
// }
