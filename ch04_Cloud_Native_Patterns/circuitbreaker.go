package ch04

import (
	"context"
	"errors"
	"sync"
	"time"
)

//La lógica que ejecutaremos. Es una función que tiene por argumento un contexto, y retorna un error y el error
type Circuit func(context.Context) (string, error)

//Abrimos el cicuito cuando el número de fallos consecutivos supere el threshold
func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	//Usaremos capture para utilziar estos datos en la función, Circuit, que devolveremos
	// Number of failures after the first
	var consecutiveFailures int = 0
	// Time of the last interaction with the downstream service
	var lastAttempt = time.Now()
	//Mutex de lectura/escritura
	var m sync.RWMutex

	// Construct and return the Circuit closure
	return func(ctx context.Context) (string, error) {
		//Comprobamos si el breaker esta abierto o no
		//El breaker se abre cuando el número de fallos consecutivos supera el threshold

		//Bloqueamos para lectura, no estamos actualizando ninguno de los datos "comunes" entre ejecuciones, los datos capturados
		m.RLock()

		//Comprobamos si el número de fallos consecutivos supera el threshold
		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			//A medida que d aumenta, los reintentos los espaciamos en el tiempo. El primer intento a los 2 segundos, el segundo a los cuatro, el tercero a los ocho...
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			//Si todavía no hemos llegado a la hora a la que tenemos que reintentar de nuevo, retorna un error
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock()

		//Si el breaker esta cerrado, es decir, no hemos superado el umbral, ejecutamos
		//Ejecutamos el circuito
		response, err := circuit(ctx) // Issue request proper

		//Bloqueamos para escritura
		m.Lock() // Lock around shared resources
		defer m.Unlock()

		//Actualizamos el timestamp de la ejecución
		lastAttempt = time.Now() // Record time of attempt
		//Si hay un error
		if err != nil { // Circuit returned an error,
			//Incrementamos el contador de fallos
			consecutiveFailures++ // so we count the failure
			//y respondemos
			return response, err // and return
		}

		//Reseteamos el contador de fallos
		consecutiveFailures = 0 // Reset failures counter
		//y respondemos
		return response, nil
	}
}
