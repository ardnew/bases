package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalString(t *testing.T) {
	for _, tc := range []struct {
		name string
		string
		rpn string
	}{
		// TODO: Add test cases.
		{
			name:   "literal",
			string: `"hello"`,
			rpn:    `"hello"`,
		},
		{
			name:   "ident",
			string: `foo`,
			rpn:    `foo`,
		},
		{
			name:   "prefix",
			string: `--(++2)`,
			rpn:    `(--, (++, 2))`,
		},
		{
			name:   "postfix",
			string: `57--`,
			rpn:    `(--, 57)`,
		},
		{
			name:   "infix",
			string: `foo * 12.1`,
			rpn:    `(*, foo, 12.1)`,
		},
		{
			name:   "order",
			string: `foo - 12.1 * 123 + 'r' / "bar" || true`,
			rpn:    `(||, (-, foo, (+, (*, 12.1, 123), (/, 'r', "bar"))), true)`,
		},
		//{
		//name: "wtfix",
		//string:    `1 + 3.14 + a-- * b & x[tt%7] | ^32 &^ ++a = ~a >> !foo - x`,
		//rpn: ,
		//},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.rpn, nil)
			// tc.assertion(t, EvalString(tc.string), fmt.Sprintf("EvalString(%v)", tc.string))
		})
	}
}
