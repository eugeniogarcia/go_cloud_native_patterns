package main

type frog struct{}

func (f frog) Says() string {
	return "ribbit!"
}

// Animal is exported as a symbol.
var Animal frog
