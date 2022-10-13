package list

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func cmpSingly(a, b *singly) bool {
	if a == nil && b == nil {
		return true // both pointers are nil
	}
	if a == nil || b == nil {
		return false // one of the pointers is nil
	}
	// both pointers are non-nil, compare their content
	return a.Node.String() == b.Node.String() &&
		cmpSingly(a.next, b.next)
}

func cmpStack(a, b *Stack) bool {
	if a == nil && b == nil {
		return true // both pointers are nil
	}
	if a == nil || b == nil {
		return false // one of the pointers is nil
	}
	// both pointers are non-nil, compare their content
	return a.len == b.len &&
		cmpSingly(a.top, b.top)
}

func comparer() cmp.Option { return cmp.Comparer(cmpStack) }

type strNode string

func (s strNode) String() string     { return string(s) }
func (s strNode) Equals(a Node) bool { return s.String() == a.String() }

func Test_Stack_Push(t *testing.T) {
	for name, tt := range map[string]struct {
		s  *Stack
		n  Node
		sw *Stack
		o  cmp.Option
	}{
		"empty": {
			&Stack{},
			strNode("abc"),
			&Stack{&singly{nil, strNode("abc")}, 1},
			comparer(),
		},
		"non-empty": {
			&Stack{&singly{&singly{nil, strNode("abc")}, strNode("def")}, 2},
			strNode("xyz"),
			&Stack{&singly{&singly{&singly{nil, strNode("abc")}, strNode("def")}, strNode("xyz")}, 3},
			comparer(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Logf("%s <- %s:", tt.s.String(), tt.n)
			tt.s.Push(tt.n)
			t.Logf("   got: %s", tt.s.String())
			t.Logf("  want: %s", tt.sw.String())
			if diff := cmp.Diff(tt.sw, tt.s, tt.o); len(diff) > 0 {
				t.Errorf("diff (-want +got):%s\n", diff)
			}
		})
	}
}

func Test_Stack_Pop(t *testing.T) {
	for name, tt := range map[string]struct {
		s  *Stack
		nw Node
		sw *Stack
		o  cmp.Option
	}{
		"empty": {
			&Stack{},
			nil,
			&Stack{},
			comparer(),
		},
		"single": {
			&Stack{&singly{nil, strNode("abc")}, 1},
			strNode("abc"),
			&Stack{},
			comparer(),
		},
		"several": {
			&Stack{&singly{&singly{&singly{nil, strNode("abc")}, strNode("def")}, strNode("xyz")}, 3},
			strNode("xyz"),
			&Stack{&singly{&singly{nil, strNode("abc")}, strNode("def")}, 2},
			comparer(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Logf("%s:", tt.s.String())
			ng := tt.s.Pop()
			t.Logf("   got: %s -> [%s]", tt.s.String(), ng)
			t.Logf("  want: %s -> [%s]", tt.sw.String(), tt.nw)
			if diff := cmp.Diff(tt.nw, ng); len(diff) > 0 {
				t.Errorf("diff (-want +got):%s\n", diff)
			}
			if diff := cmp.Diff(tt.sw, tt.s, tt.o); len(diff) > 0 {
				t.Errorf("diff (-want +got):%s\n", diff)
			}
		})
	}
}
