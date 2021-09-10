package main

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Host string
	Port uint
	Tags map[string]string
}

//Uso de field tags para personalizar el marshalling/unmarshaling a json
type Tagged struct {
	// CustomKey will appear in JSON as the key "custom_key".
	CustomKey string `json:"custom_key"`
	// OmitEmpty will appear in JSON as "OmitEmpty" (the default),
	// but will only be written if it contains a nonzero value.
	OmitEmpty string `json:",omitempty"`
	// IgnoredName will always be ignored.
	IgnoredName string `json:"-"`
	// TwoThings will appear in JSON as the key "two_things",
	// but only if it isn't empty.
	TwoThings string `json:"two_things,omitempty"`
}

func main() {

	c := Config{
		Host: "localhost",
		Port: 1313,
		Tags: map[string]string{"env": "dev"},
	}

	//Crea Json
	bytes, _ := json.Marshal(c)
	fmt.Println("Convierte un objeto a JSON")
	fmt.Println(string(bytes))
	fmt.Println("")

	//Crea Json con formato
	//Los argumento son el prefijo, y que caracter se usara cuando haya que poner un espacio
	fmt.Println("Convierte un objeto a JSON, con formato")
	bytes, _ = json.MarshalIndent(c, "prefijo", "_")
	fmt.Println(string(bytes))
	fmt.Println("")
	bytes, _ = json.MarshalIndent(c, "", " ")
	fmt.Println(string(bytes))
	fmt.Println("")

	//Convierte el JSon en un objeto
	d := Config{}
	fmt.Println("Convierte un JSON a un objeto")
	_ = json.Unmarshal(bytes, &d)
	fmt.Println(d)
	fmt.Println("")

	//Convierte el JSon en un objeto
	fmt.Println("Convierte un JSON a un objeto generico")
	var e interface{}
	_ = json.Unmarshal(bytes, &e)

	//Es preciso hacer el cast antes de usarlo. El unmarshal convierte la información a un mapa de interface{}
	m := e.(map[string]interface{})
	fmt.Printf("<%T> %v\n", m, m)
	fmt.Printf("<%T> %v\n", m["Foo"], m["Foo"])
	fmt.Printf("<%T> %v\n", m["Number"], m["Number"])
	fmt.Printf("<%T> %v\n", m["Tags"], m["Tags"])
	fmt.Println("")

	personaliza1 := Tagged{
		CustomKey:   "clave",
		IgnoredName: "ignorame",
		OmitEmpty:   "no esta vacio",
		TwoThings:   "dos cosas",
	}
	personaliza2 := Tagged{
		CustomKey:   "",
		IgnoredName: "",
		OmitEmpty:   "",
		TwoThings:   "",
	}

	//Crea Json usando tags
	fmt.Println("Convierte un objeto a JSON usando tags")
	bytes_per1, _ := json.MarshalIndent(personaliza1, "", " ")
	fmt.Println(string(bytes_per1))
	fmt.Println("")
	bytes_per2, _ := json.MarshalIndent(personaliza2, "", " ")
	fmt.Println(string(bytes_per2))
	fmt.Println("")

	//Convierte el JSon en un objeto usando tags
	resp_per1 := Tagged{}
	resp_per2 := Tagged{}
	fmt.Println("Convierte un JSON a un objeto usando tags")
	_ = json.Unmarshal(bytes_per1, &resp_per1)
	_ = json.Unmarshal(bytes_per2, &resp_per2)
	fmt.Printf("CustomKey=%s, TwoThings=%s, OmitEmpty=%s\n", resp_per1.CustomKey, resp_per1.TwoThings, resp_per1.OmitEmpty)
	fmt.Printf("CustomKey=%s, TwoThings=%s, OmitEmpty=%s\n", resp_per2.CustomKey, resp_per2.TwoThings, resp_per2.OmitEmpty)
	fmt.Println("")

}
