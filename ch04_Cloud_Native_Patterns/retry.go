package ch04

import (
	"context"
	"log"
	"time"
)

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		//Loop infinito
		for r := 0; ; r++ {
			//Hacemos la llamada a la lógica
			response, err := effector(ctx)
			//Si no hay error, o se superaron los intentos, retornamos la respuesta de la lógica...
			if err == nil || r >= retries {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)
			//Esperamos delay segundos hasta intentarlo de nuevo
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
