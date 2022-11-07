package list

import "strings"

// singly represents a singly-linked list composed of Node elements.
type singly struct {
	next *singly
	Node
}

// String returns a comprehensive string representation of the list.
func (n *singly) String() string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(n.Node.String())
	b.WriteRune(',')
	b.WriteString(n.next.String())
	return b.String()
}
