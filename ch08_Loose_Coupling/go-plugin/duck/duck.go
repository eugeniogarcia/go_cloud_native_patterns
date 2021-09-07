package main

type duck struct{}

func (d duck) Says() string {
	return "quack!"
}

// Animal is exported as a symbol.
var Animal duck
