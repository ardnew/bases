package eval

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalString(t *testing.T) {
	for name, tc := range map[string]struct {
		assertion assert.ErrorAssertionFunc
		string
	}{
		// TODO: Add test cases.
		// "literal": {
		// 	string:    `"hello"`,
		// 	assertion: assert.NoError,
		// },
		// "ident": {
		// 	string:    "foo",
		// 	assertion: assert.NoError,
		// },
		"prefix": {
			string:    "--(++2)",
			assertion: assert.NoError,
		},
		// "postfix": {
		// 	string:    "57--",
		// 	assertion: assert.NoError,
		// },
		// "infix": {
		// 	string:    "foo * 12.1",
		// 	assertion: assert.NoError,
		// },
		// "order": {
		// 	string:    "foo * 12.1 + 123 - 'r'",
		// 	assertion: assert.NoError,
		// },
		// "wtfix": {
		// 	string:    "1 + 3.14 + a-- * b & x[tt%7] | ^32 &^ ++a = ~a >> !foo - x",
		// 	assertion: assert.NoError,
		// },
	} {
		t.Run(name, func(t *testing.T) {
			tc.assertion(t, EvalString(tc.string), fmt.Sprintf("EvalString(%v)", tc.string))
		})
	}
}
