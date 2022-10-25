package op

// Go lexer currently recognizes 89 tokens — not just operators — but we need
// to have enough indices for an operator token constant defined anywhere in the
// enumerated list of tokens. This enables constant-time lookup for operator
// tokens, instead of a map's relatively more-expensive hashing function.
const maxOperators = 128 // min({ x in 2**N | x > 89 })
