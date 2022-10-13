package expr

import "github.com/ardnew/bases/io/syntax/op"

type Expr interface{}

type Ident struct {
	ident string
}

type Assign struct {
	expr  Expr
	ident string
}

type Infix struct {
	left, right Expr
	op          op.Binary
}

type Prefix struct {
	expr Expr
	op   op.Unary
}

type Postfix struct {
	expr Expr
	op   op.Unary
}

type Mixfix struct {
	cond, left, right Expr
}
