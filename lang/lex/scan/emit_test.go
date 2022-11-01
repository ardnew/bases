package scan

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ardnew/bases/lang/lex/sym"
)

func TestMake(t *testing.T) {
	for name, tt := range map[string]struct {
		args struct {
			buff []byte
		}
		wantEmit Emit
	}{
		// TODO: Add test cases.
		"1": {
			args: struct{ buff []byte }{
				buff: []byte("x + 3 -- foo"),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			em := Make(tt.args.buff)
			ch := make(sym.Stream)
			go func(e Emit, c sym.Stream) {
				for e != nil {
					e = e(c)
				}
				close(c)
			}(em, ch)
			for c := range ch {
				fmt.Println(c)
			}
		})
	}
}

func TestEmit_Undo(t *testing.T) {
	for name, tt := range map[string]struct {
		e    Emit
		args struct {
			s sym.Symbol
		}
		want Emit
	}{
		// TODO: Add test cases.
		"": {},
	} {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.e.Undo(tt.args.s)); len(diff) > 0 {
				t.Errorf("diff (-want +got):%s\n", diff)
			}
		})
	}
}
