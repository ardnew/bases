package list

// Exported system architecture limits
const (
	UINT_MAX = ^uint(0)
	INT_MAX  = UINT_MAX >> 1
)

// Unexported logical constraints
const (
	maxLen = INT_MAX
)
