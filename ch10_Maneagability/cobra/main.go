package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var strp string
var intp int
var boolp bool

//Comando base
var rootCmd = &cobra.Command{
	Use:  "cng (nombre corto)",
	Long: "A super simple command. (descripción larga)",
}

//Crea un comando llamado Flags
var flagsCmd = &cobra.Command{
	Use:   "banderas",
	Short: "Experimenta con flags (descripción corta)",
	Long:  "Muestra el valor de los flags. (descripción larga)",
	Run:   flagsFunc,
}

func flagsFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Fstring:", strp)
	fmt.Println("Finteger:", intp)
	fmt.Println("Fboolean:", boolp)
	fmt.Println("args:", args)
}

func init() {
	//Definimos tres flags
	flagsCmd.Flags().StringVarP(&strp, "Fstring", "s", "foo", "a string")
	flagsCmd.Flags().IntVarP(&intp, "Fnumber", "n", 42, "an integer")
	flagsCmd.Flags().BoolVarP(&boolp, "Fboolean", "b", false, "a boolean")

	//Y los añadimos al comando
	rootCmd.AddCommand(flagsCmd)

	helloWorldCmd.Flags().BoolVarP(&mayuscula, "mayuscula", "m", false, "si mostramos el texto en mayusculas o no")
	//Añadimos el segundo comando
	rootCmd.AddCommand(helloWorldCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
