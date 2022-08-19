package main

import "github.com/ardnew/bases/ui"

func main() {
	ui := ui.New()
	if err := ui.Run(); err != nil {
		panic(err)
	}
}
