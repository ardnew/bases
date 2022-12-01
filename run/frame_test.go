package run

import (
	"fmt"
	"runtime"
	"testing"
)

func TestCallstack(t *testing.T) {
	dump := func(f []runtime.Frame) {
		fmt.Printf("%d Frames:\n", len(f))
		for i, c := range f {
			fmt.Printf("  %02d: %+v\n", i, c)
		}
		fmt.Println()
	}
	tests := []struct {
		name string
		size int
	}{
		{"size=0", 0},
		{"size=1", 1},
		{"size=2", 2},
		{"size=3", 3},
		{"size=4", 4},
		{"size=5", 5},
		{"size=-0", -0},
		{"size=-1", -1},
		{"size=-2", -2},
		{"size=-3", -3},
		{"size=-4", -4},
		{"size=-5", -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			dump(Callstack(tt.size))
		})
	}
}
