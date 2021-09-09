package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var mayuscula bool

//Segundo comando
var helloWorldCmd = &cobra.Command{
	Use:   "hola",
	Short: "Imprime \"Hola, Mundo\"",
	Long:  "Este comando imprime \"Hola, Mundo\"",
	Run:   helloWorldFunc,
	Args:  cobra.MaximumNArgs(1),
}

func helloWorldFunc(cmd *cobra.Command, args []string) {
	name := "Mundo"

	if len(args) > 0 {
		name = args[0]
	}

	if mayuscula {
		fmt.Printf("HOLA, %s.\n", strings.ToUpper(name))
	} else {
		fmt.Printf("Hola, %s.\n", name)
	}
}
