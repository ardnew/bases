package lex

import (
	"testing"
)

func TestScan_All(t *testing.T) {
	for name, tt := range map[string]struct {
		s *Scan
		// want []Token
	}{
		// TODO: Add test cases.
		"simple": {
			s: NewScan(5).Init("1 + 3.14"),
			// want: NewScan().Init("1 + 3").All(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, tok := range tt.s.All() {
				t.Logf("%s %s %q", tok.Pos, tok.Tok, tok.Lit)
			}
			t.Log()
			for i, tok := range tt.s.prev {
				t.Logf("[%d] %s %s %q", i, tok.Pos, tok.Tok, tok.Lit)
			}
			//if diff := cmp.Diff(tt.want, tt.s.All()); len(diff) > 0 {
			//	t.Errorf("diff (-want +got):%s\n", diff)
			//}
		})
	}
}
