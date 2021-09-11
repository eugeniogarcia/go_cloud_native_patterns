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

func (c *Config) imprime() {
	fmt.Printf("%s:%d\n", c.Host, c.Port)
	fmt.Println(c.Tags)
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

func (t *Tagged) imprimeJSON() ([]byte, error) {
	fmt.Println("Convierte un objeto a JSON usando tags")
	bytes_per, err := json.MarshalIndent(t, "", " ")
	fmt.Println(string(bytes_per))
	fmt.Println("")
	return bytes_per, err
}

func (t *Tagged) leeJSON(bytes_per []byte) error {
	fmt.Println("Convierte un JSON a un objeto usando tags")
	err := json.Unmarshal(bytes_per, t)
	fmt.Printf("CustomKey=%s, TwoThings=%s, OmitEmpty=%s\n", t.CustomKey, t.TwoThings, t.OmitEmpty)
	fmt.Println("")

	return err
}

func (t *Tagged) iguales(a Tagged) bool {
	if t.CustomKey == a.CustomKey && t.OmitEmpty == a.OmitEmpty && t.TwoThings == a.TwoThings {
		fmt.Println("Son iguales")
		return true
	} else {
		fmt.Println("Son distintos")
		return false
	}
}
