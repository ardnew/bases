package op

// The Go lexer [go/scanner.Scanner] currently uses 90 constants (iota..89,
// including some bookkeeping constants) defined in [go/token].
//
// So in order to map these tokens to anything, e.g., parser functions, we need
// an array of length N >= 90 for perfect hashing with constant-time access.
const maxOperators = 128

// Every Operator is an n-ary function, where n <= MaxArity.
const MaxArity = 2
