package list

// Exported system architecture limits
const (
	UintMax = ^uint(0)
	IntMax  = UintMax >> 1
)

// Exported logical constraints
const (
	MaxLen = IntMax
)
