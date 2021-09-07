package main

import (
	"fmt"
	"os"
	"plugin"
)

// Sayer says what an animal says.
type Sayer interface {
	Says() string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: run main/main.go animal")
		os.Exit(1)
	}

	// Get the animal name, and build the path where we expect to
	// find the corresponding shared object (.so) file.
	name := os.Args[1]
	module := fmt.Sprintf("./%s/%s.so", name, name)

	// Does the file exist?
	_, err := os.Stat(module)
	if os.IsNotExist(err) {
		fmt.Println("can't find an animal named", name)
		os.Exit(1)
	}

	// Open our plugin. and returns a *plugin.Plugin.
	p, err := plugin.Open(module)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Lookup searches for a symbol, which can be any exported variable
	// or function, named "Animal" in plugin p.
	symbol, err := p.Lookup("Animal")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Asserts that the symbol interface holds an Sayer.
	animal, ok := symbol.(Sayer)
	if !ok {
		fmt.Println("that's not an Sayer")
		os.Exit(1)
	}

	// Now we can use our loaded plugin!
	fmt.Printf("A %s says: %q\n", name, animal.Says())
}
