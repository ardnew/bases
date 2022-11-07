package list

import (
	"sync"
)

type SyncStack struct {
	sync.RWMutex
	Stack
}

// Len returns the number of elements on the stack.
func (s *SyncStack) Len() int {
	s.RLock()
	defer s.RUnlock()
	return s.Stack.Len()
}

// Push adds a new element to the top of the stack.
//
// Nil elements are ignored and will not be added.
func (s *SyncStack) Push(a ...Node) {
	s.Lock()
	defer s.Unlock()
	s.Stack.Push(a...)
}

// Pop removes and returns the element on top of the stack.
//
// Nil is returned from an empty stack.
func (s *SyncStack) Pop() (a Node) {
	s.Lock()
	defer s.Unlock()
	return s.Stack.Pop()
}

// Peek returns the element on top of the stack.
// Unlike [Pop], it does not remove the element from the stack.
//
// Nil is returned from an empty stack.
func (s *SyncStack) Peek() (a Node) {
	s.RLock()
	defer s.RUnlock()
	return s.Stack.Peek()
}

// String returns a comprehensive string representation of the stack.
func (s *SyncStack) String() string {
	s.RLock()
	defer s.RUnlock()
	return s.Stack.String()
}
