# Cobra

Podemos encontrar el [manual aquí](https://github.com/spf13/cobra/blob/master/user_guide.md).

## How to use it

- Estructura. La estrcutura consiste en _comandos_ que pueden usar _flags_. Tendremos siempre un comando base que equivale al propio programa, nuestro ejecutable, y comandos adicionales que están asociados a una función. Estos comandos pueden usar los _flags_ que pasamos como parametros con el comando. Los pasos para implementar todo esto son: 

- Declarar las variables que usaremos como flags a nivel de paquete:

```go
var strp string
var intp int
var boolp bool
```

- Creamos la estructura que representa la base de Cobra.

```go
//Comando base
var rootCmd = &cobra.Command{
	Use:  "cng (nombre del ejecutable)",
	Long: "A super simple command. (descripción larga)",
}
```

- Creamos un comando llamado `banderas`. _Run_ es la función que implementa lo que se tiene que hacer cuando se ejecuta:

```go
//Crea un comando llamado Flags
var flagsCmd = &cobra.Command{
	Use:   "banderas",
	Short: "Experimenta con flags (descripción corta)",
	Long:  "Muestra el valor de los flags. (descripción larga)",
	Run:   flagsFunc,
}
```

- La función que implementa el comando. En este caso simplemente imprimimos el valor de los flags y los argumentos que se pasan al programa:

```go
func flagsFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Fstring:", strp)
	fmt.Println("Finteger:", intp)
	fmt.Println("Fboolean:", boolp)
	fmt.Println("args:", args)
}
```

- Podemos definir varios comandos. Por ejemplo aquí tenemos otro llamado _hola_. Indicamos que el número máximo de argumentos que podemos pasar al programa cuando usemos este comando es uno:

```go
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
```

- Finalmente usamos la función _init_ para añadir los flags a cada comando, y los comandos al comando base:

```go
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
```

## Ejemplos

Podemos pedir ayuda usando _-h_ o _--help_:

```ps
go run . -h

A super simple command. (descripción larga)

Usage:
  cng [command]

Available Commands:
  banderas    Experimenta con flags (descripción corta)
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  hola        Imprime "Hola, Mundo"

Flags:
  -h, --help   help for cng

Use "cng [command] --help" for more information about a command.
```

Usamos el comando _banderas_ sin pasar flags, y con un argumento:

```ps
go run . banderas "arg1"

Fstring: foo
Finteger: 42
Fboolean: false
args: [arg1]
```

Usamos el comando _hola_ con un argumento:

```ps
go run . hola           

Hola, Mundo.
```

Usamos el comando _hola_ con un argumento, e indicando que se usen mayusculas:

```ps
go run . hola -m=true "eugenio"

HOLA, EUGENIO.
```

Vemos la ayuda del comando _hola_:

```ps
go run . hola -h               

Este comando imprime "Hola, Mundo"

Usage:
  cng hola [flags]

Flags:
  -h, --help        help for hola
  -m, --mayuscula   si mostramos el texto en mayusculas o no
```