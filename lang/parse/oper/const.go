package oper

// The Go lexer [go/scanner.Scanner] currently uses 90 constants (iota..89,
// including some bookkeeping constants) defined in [go/token].
//
// A perfect hashing function with constant-time access can trivially be
// implemented with an array of length N >= 90.
//
// Using 128 because its a power of 2 (perhaps neatly-aligned in memory), and
// it leaves sufficient space (but not excessive) for updates to [go/token].
const MaxOperators = 128

// Every Operator has no more than N operands, N <= MaxArity.
//
// In other words, every Operator is an N-ary function.
const MaxArity = 2
