package main

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru" //Cache lru de hashicorp
)

var cache *lru.Cache

//Crea una cache
func init() {
	cache, _ = lru.NewWithEvict(2,
		func(key interface{}, value interface{}) {
			fmt.Printf("Evicted: key=%v value=%v\n", key, value)
		},
	)
}
func main() {
	//Añade dos entradas
	cache.Add(1, "a") // adds 1
	cache.Add(2, "b") // adds 2; cache is now at capacity
	if valor, ok := cache.Get(1); ok {
		fmt.Printf("Valor de la clave 1: %s\n", valor.(string))
	} // "a true"; 1 now most recently used
	cache.Remove(2)
	if !cache.Contains(2) {
		fmt.Printf("La clave 2 ya no está en la cache. Hay %d entradas en la cache\n", cache.Len())
	}
	cache.Add(3, "c")
	//Como el tamaño de la cache es 2, se provoca la salida de un elemento de la cache
	cache.Add(4, "d")
}
