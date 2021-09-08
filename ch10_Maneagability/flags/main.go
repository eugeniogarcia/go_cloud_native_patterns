package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Declare a string flag with a default value "foo"
	// and a short description. It returns a string pointer.
	strp := flag.String("Fstring", "foo", "a string")

	// Declare number and boolean flags, similar to the string flag.
	intp := flag.Int("Fnumber", 42, "an integer")
	boolp := flag.Bool("Fboolean", false, "a boolean")

	// Call flag.Parse() to execute command-line parsing.
	flag.Parse()

	// Print the parsed options and trailing positional arguments.
	fmt.Println("string:", *strp)
	fmt.Println("integer:", *intp)
	fmt.Println("boolean:", *boolp)
	fmt.Println("args:", flag.Args())

	//Configuración con variables de entorno. Recuperamos el valor
	name := os.Getenv("NAME")
	place := os.Getenv("CITY")

	fmt.Printf("%s lives in %s.\n", name, place)

	//Con lookup podemos distinguir si la variable no esta informada del caso en el que esta informada con blanco
	key := "NAME"
	if val, ok := os.LookupEnv(key); ok {
		fmt.Printf("%s=%s\n", key, val)
	} else {
		fmt.Printf("%s not set\n", key)
	}
}
