package ch04

import (
	"context"
	"sync"
	"time"
)

type Future interface {
	Result() (string, error)
}

type InnerFuture struct {
	once sync.Once
	wg   sync.WaitGroup

	res   string
	err   error
	resCh <-chan string
	errCh <-chan error
}

func (f *InnerFuture) Result() (string, error) {
	//Se ejecuta una vez, y se bloquea la ejecución hasta no recibir algo por los canales
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	//Bloquea hasta que once termine
	f.wg.Wait()

	return f.res, f.err
}

//Lógica "lenta". Devuelbe un Future
func SlowFunction(ctx context.Context) Future {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <-time.After(time.Second * 2):
			resCh <- "I slept for 2 seconds"
			errCh <- nil
		case <-ctx.Done():
			resCh <- ""
			errCh <- ctx.Err()
		}
	}()

	//Crea el future
	return &InnerFuture{resCh: resCh, errCh: errCh}
}