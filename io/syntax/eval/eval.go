package eval

import (
	"fmt"
	"time"

	"github.com/ardnew/bases/io/syntax/parse"
	"github.com/ardnew/bases/io/syntax/parse/expr"
	"github.com/ardnew/bases/io/syntax/parse/lex/scan"
)

type eval struct {
	parse.State
	*scan.Scan
}

func EvalString(s string) error {
	e := eval{
		State: expr.Parse,
		Scan:  scan.NewScan().Init([]byte(s)),
	}
	return e.Parse()
}

func (e *eval) Parse() error {
	var f, hello, world func()
	hello = func() {
		fmt.Print("hello ")
		f = world
	}
	world = func() {
		fmt.Println("world!")
		f = hello
	}
	f = hello
	for f != nil {
		f()
		time.Sleep(time.Second)
	}

	for e.State != nil {
		e.State = e.State(e)
	}
	return e
}
