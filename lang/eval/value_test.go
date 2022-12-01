package eval

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalString(t *testing.T) {
	for _, tc := range []struct {
		name string
		string
		sex string
	}{
		{
			name:   "literal",
			string: `"foo"`,
			sex:    `"foo"`,
		},

		{
			name:   "identifier",
			string: `bar`,
			sex:    `bar`,
		},

		{
			name:   "prefix A",
			string: `++a`,
			sex:    `(++ a)`,
		},

		{
			name:   "infix A",
			string: `1 + 2 * 3`,
			sex:    `(+ 1 (* 2 3))`,
		},

		{
			name:   "infix B",
			string: `a + b * c * d + e`,
			sex:    `(+ (+ a (* (* b c) d)) e)`,
		},

		// TODO: Add test cases.
		// {
		// 	name:   "literal",
		// 	string: `"hello"`,
		// 	sex:    `"hello"`,
		// },
		// {
		// 	name:   "ident",
		// 	string: `foo`,
		// 	sex:    `foo`,
		// },
		// {
		// 	name:   "prefix",
		// 	string: `--(++2)`,
		// 	sex:    `(--, (++, 2))`,
		// },
		// {
		// 	name:   "postfix",
		// 	string: `57--`,
		// 	sex:    `(--, 57)`,
		// },
		// {
		// 	name:   "infix",
		// 	string: `foo * 12.1`,
		// 	sex:    `(*, foo, 12.1)`,
		// },
		// {
		// 	name:   "order",
		// 	string: `foo - 12.1 * 123 + 'r' / "bar" || true`,
		// 	sex:    `(||, (-, foo, (+, (*, 12.1, 123), (/, 'r', "bar"))), true)`,
		// },
		//{
		//name: "wtfix",
		//string:    `1 + 3.14 + a-- * b & x[tt%7] | ^32 &^ ++a = ~a >> !foo - x`,
		//sex: ,
		//},
	} {
		t.Run(tc.name, func(t *testing.T) {
			val := New()
			io.Copy(val, strings.NewReader(tc.string))
			assert.Equal(t, tc.sex, val.String())
		})
	}
}
