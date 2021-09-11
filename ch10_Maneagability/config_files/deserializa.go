package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-yaml/yaml"
)

func deserializaJSON(bytes []byte) (Config, error) {
	//Convierte el JSon en un objeto
	d := Config{}
	fmt.Println("Convierte un JSON a un objeto")
	err := json.Unmarshal(bytes, &d)
	fmt.Println(d)
	fmt.Println("")

	//Convierte el JSon en un objeto
	fmt.Println("Convierte un JSON a un objeto generico")
	var e interface{}
	err = json.Unmarshal(bytes, &e)

	//Es preciso hacer el cast antes de usarlo. El unmarshal convierte la información a un mapa de interface{}
	m := e.(map[string]interface{})
	fmt.Printf("<%T> %v\n", m, m)
	fmt.Printf("<%T> %v\n", m["Host"], m["Host"])
	fmt.Printf("<%T> %v\n", m["Port"], m["Port"])
	fmt.Printf("<%T> %v\n", m["Tags"], m["Tags"])
	fmt.Println("")

	return d, err
}

func deserializaYaml(bytes_yaml []byte) (Config, error) {
	//Convierte el JSon en un objeto
	d := Config{}
	fmt.Println("Convierte un Yaml a un objeto")
	err := yaml.Unmarshal(bytes_yaml, &d)
	fmt.Println(d)
	fmt.Println("")

	//Convierte el Yaml en un objeto
	var e interface{}
	fmt.Println("Convierte un Yaml a un objeto generico")
	err = yaml.Unmarshal(bytes_yaml, &e)

	//Es preciso hacer el cast antes de usarlo.	 OJO QUE EL TIPO ES DIFERENTE AL QUE OBTENIAMOS CON LOS JSONS
	//OTRA DIFERENCIA ES QUE LOS KEYS SE PONEN EN MINUSCULA
	m_yaml := e.(map[interface{}]interface{})
	fmt.Printf("<%T> %v\n", m_yaml, m_yaml)
	fmt.Printf("<%T> %v\n", m_yaml["host"], m_yaml["host"])
	fmt.Printf("<%T> %v\n", m_yaml["port"], m_yaml["port"])
	fmt.Printf("<%T> %v\n", m_yaml["tags"], m_yaml["tags"])
	fmt.Println("")

	return d, err
}
