package ch04

import (
	"context"
	"sync"
	"time"
)

func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time = time.Now()
	var ticker *time.Ticker
	// para controlar que se llame solo una vez
	var once sync.Once
	var m sync.Mutex
	// cachear la respuesta
	var result string
	var err error

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()

		//Calcula hasta cuando tenemos que esperar para hacer la llamada al circuito
		threshold = time.Now().Add(d)

		//Esto lo hace una sola vez
		once.Do(func() {
			//Arranca el temporizador, 100 milisegundos
			ticker = time.NewTicker(time.Millisecond * 100)

			go func() {
				defer func() {
					m.Lock()
					//Para el temporizador
					ticker.Stop()
					//Reseteamos once, para que se pueda hacer otra llamada al circuito
					once = sync.Once{}
					m.Unlock()
				}()

				for {
					select {
					//Cada 100ms comprobamos
					case <-ticker.C:
						m.Lock()
						//Si hemos superado el threshold
						if time.Now().After(threshold) {
							//si lo hemos hecho llamamos al circuito
							result, err = circuit(ctx)
							m.Unlock()
							return
						}
						m.Unlock()
					case <-ctx.Done():
						m.Lock()
						result, err = "", ctx.Err()
						m.Unlock()
						return
					}
				}
			}()
		})

		return result, err
	}
}
