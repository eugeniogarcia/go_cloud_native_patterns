package ch04

import (
	"sync"
)

func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int) // The shared output channel

	var wg sync.WaitGroup // Used to automatically close dest
	// when all sources are closed
	wg.Add(len(sources))         // Set size of the WaitGroup
	for _, ch := range sources { // Start a goroutine for each source
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				dest <- n
			}
		}(ch)
	}
	go func() { // Start a goroutine to close dest
		wg.Wait() // after all sources close
		close(dest)
	}()

	return dest
}
