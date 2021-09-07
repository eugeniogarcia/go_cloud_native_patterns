package main

type fox struct{}

func (f fox) Says() string {
	return "ring-ding-ding-ding-dingeringeding!"
}

// Animal is exported as a symbol.
var Animal fox
