package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	c := Config{
		Host: "localhost",
		Port: 1313,
		Tags: map[string]string{"env": "dev"},
	}

	//Serializa
	bytes, _ := serializaJSON(c)
	bytes_yaml, _ := serializaYaml(c)

	//Deserializa
	d, _ := deserializaJSON(bytes)
	d, _ = deserializaYaml(bytes_yaml)

	//Compara
	if d.Host == c.Host && d.Port == c.Port && c.Tags["env"] == d.Tags["env"] {
		fmt.Println("Son iguales!!")
	} else {
		panic("No son iguales")
	}

	//Prueba con tagged elements
	personaliza := Tagged{
		CustomKey:   "clave",
		IgnoredName: "ignorame",
		OmitEmpty:   "no esta vacio",
		TwoThings:   "dos cosas",
	}
	//Serializa
	bytes_per, _ := personaliza.imprimeJSON()

	//Deserializa
	resp_per := Tagged{}
	resp_per.leeJSON(bytes_per)

	//Compara
	if !resp_per.iguales(personaliza) {
		panic("distintos")
	}

	personaliza = Tagged{
		CustomKey:   "",
		IgnoredName: "",
		OmitEmpty:   "",
		TwoThings:   "",
	}

	//Serializa
	bytes_per, _ = personaliza.imprimeJSON()

	//Deserializa
	resp_per = Tagged{}
	resp_per.leeJSON(bytes_per)

	//Compara
	if !resp_per.iguales(personaliza) {
		panic("distintos")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	//Lee la configuración de un archivo
	go func(duracion int) {
		defer wg.Done()

		ticker := time.NewTicker(time.Second * 2)
		cuenta_atras := duracion
		for range ticker.C {
			if cuenta_atras%2 == 0 {
				//Mostramos la configuración actual
				config.imprime()
			}

			cuenta_atras--
			if cuenta_atras < 0 {
				break
			}
			if cuenta_atras == 5 {
				//Cancelamos el monitoreo de la configuración
				cancela()
				println("No monitorizamos los cambios en el archivo de configuracion")
			}
		}
	}(10)

	wg.Wait()
}
