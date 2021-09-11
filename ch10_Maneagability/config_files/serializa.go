package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-yaml/yaml"
)

func serializaJSON(c Config) ([]byte, error) {
	//Crea Json
	bytes, err := json.Marshal(c)
	fmt.Println("Convierte un objeto a JSON")
	fmt.Println(string(bytes))
	fmt.Println("")

	//Crea Json con formato
	//Los argumento son el prefijo, y que caracter se usara cuando haya que poner un espacio
	fmt.Println("Convierte un objeto a JSON, con formato")
	bytes, err = json.MarshalIndent(c, "prefijo", "_")
	fmt.Println(string(bytes))
	fmt.Println("")
	bytes, err = json.MarshalIndent(c, "", " ")
	fmt.Println(string(bytes))
	fmt.Println("")

	return bytes, err
}

func serializaYaml(c Config) ([]byte, error) {
	bytes_yaml, err := yaml.Marshal(c)
	fmt.Println("Convierte un objeto a yaml")
	fmt.Println(string(bytes_yaml))
	fmt.Println("")
	return bytes_yaml, err
}
