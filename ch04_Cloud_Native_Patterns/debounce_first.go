package ch04

import (
	"context"
	"sync"
	"time"
)

//Hace una llamada a un circuito. Hasta que no  y en los proximos d segundos no se hace una nueva llamada, sino que se reutiliza la respuesta cacheada
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	//Datos a capturar por el circuito
	//Para controlar el tiempo hasta permitir una nueva llamada
	var threshold time.Time
	var m sync.Mutex
	//Para cachear la respuesta
	var result string
	var err error

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			//Se actualiza el threshold
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		//Si no se ha superado el threshold, usamos la información cacheada
		if time.Now().Before(threshold) {
			return result, err
		}

		//Llamamos al circuito, guardamos la respuesta
		result, err = circuit(ctx)
		//y la retornamos
		return result, err
	}
}
