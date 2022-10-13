package lex

import "go/token"

type Token struct {
	Pos token.Position
	Tok token.Token
	Lit string
}
