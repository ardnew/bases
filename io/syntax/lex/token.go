package lex

import "go/token"

type Token struct {
	Lit string
	Pos token.Position
	Tok token.Token
}
