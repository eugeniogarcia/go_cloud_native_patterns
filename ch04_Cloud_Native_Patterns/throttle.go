package ch04

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	//tokens. Indica la cuota de ejecuciones disponible
	var tokens = max
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		//Se ejecuta una vez
		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				//Se ejecuta de forma indefinada hasta que el contexto se de por cerrado
				for {
					select {
					case <-ctx.Done():
						return

					//Pasado este tiempo...
					case <-ticker.C:
						m.Lock()
						//rellenamos la cuota de tokens...
						t := tokens + refill
						//...asegurandonos de que no se supera la cuota máxima
						if t > max {
							t = max
						}
						tokens = t
						m.Unlock()
					}
				}
			}()
		})

		m.Lock()
		defer m.Unlock()

		//Comprobamos si la cuota de tokens se ha agotado
		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}

		//Consumimos un token
		tokens--

		//Sino se ha agotado, se ejecuta la lógica
		return e(ctx)
	}
}
