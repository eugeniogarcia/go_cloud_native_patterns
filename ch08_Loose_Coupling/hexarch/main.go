package main

import (
	"log"

	"egsmartin.com/hexarch/core"
	"egsmartin.com/hexarch/frontend"
	"egsmartin.com/hexarch/transact"
)

func main() {
	//Adaptador en la arquitectura hexagonal
	//Usamos esta factoria para elegir el transaction logger
	tl, _ := transact.NewTransactionLogger("file")

	//Creamos la funcionalidad. Usaremos el transaction logger que acabamos de crear. El transaction logger es el componente que hace de Adapter en la arquitectura hexagonal
	store := core.NewKeyValueStore().WithTransactionLogger(tl)
	store.Restore()

	//Port en la arquitectura hexagonal
	//Elegimos el frontend, o dicho de otra forma, como vamos a exponer la funcionalidad
	fe, _ := frontend.NewFrontEnd("rest")

	//Exponemos la logica core, con los puertos y adaptadores que hemos creado con las factorias
	log.Fatal(fe.Start(store))
}
