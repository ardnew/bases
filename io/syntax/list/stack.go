package list

import (
	"strconv"
	"strings"
)

// Stack represents a last-in, first-out (LIFO) data structure composed of
// [Node] elements. The zero value is ready to be used.
type Stack struct {
	top *singly
	len uint
}

// Len returns the number of elements on the stack.
func (s *Stack) Len() int {
	// The comparison to [maxLen] in [Push] guarantees this will fit.
	return int(s.len)
}

// Push adds a new element to the top of the stack.
//
// Nil elements are ignored and will not be added.
func (s *Stack) Push(a Node) {
	if a != nil && s.len < maxLen {
		s.top = &singly{Node: a, next: s.top}
		s.len++
	}
}

// Pop removes and returns the element on top of the stack.
//
// Nil is returned from an empty stack.
func (s *Stack) Pop() (a Node) {
	if s.len > 0 {
		a = s.top.Node
		s.top = s.top.next
		s.len--
	}
	return
}

// Peek returns the element on top of the stack.
// Unlike [Pop], it does not remove the element from the stack.
//
// Nil is returned from an empty stack.
func (s *Stack) Peek() (a Node) {
	if s.len > 0 {
		return s.top.Node
	}
	return
}

// String returns a comprehensive string representation of the stack.
func (s *Stack) String() string {
	var b strings.Builder
	b.WriteRune('{')
	b.WriteString(strconv.FormatUint(uint64(s.len), 10))
	b.WriteString(": ")
	b.WriteString(s.top.String())
	b.WriteRune('}')
	return b.String()
}
