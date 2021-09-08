package main

import (
	"context"
	"time"
)

// Effector is the function that you want to subject to throttling.
type Effector func(context.Context) (string, error)

// Throttled wraps an Effector. It accepts the same parameters, plus a
// "uid" string that represents a caller identity. It returns the same,
// plus a bool that's true if the call is not throttled.
type Throttled func(context.Context, string) (bool, string, error)

// A bucket tracks the requests associated with a uid. Tendremos un token por usuario
type bucket struct {
	tokens uint
	time   time.Time
}

//No es production ready. No gestionamos la concurrencia, ni la posibilidad de que un usuario deje de llamar - eliminar al usuario del mapa buckets

// Throttle accepts an Effector function, and returns a Throttled
// function with a per-uid token bucket with a capacity of max
// that refills at a rate of refill tokens every d.
func Throttle(e Effector, max uint, refill uint, d time.Duration) Throttled {
	// buckets maps uids to specific buckets
	buckets := map[string]*bucket{}

	return func(ctx context.Context, uid string) (bool, string, error) {
		b := buckets[uid]

		// This is a new entry! It passes. Assumes that capacity >= 1.
		if b == nil {
			buckets[uid] = &bucket{tokens: max - 1, time: time.Now()}

			str, err := e(ctx)
			return true, str, err
		}

		// Calculate how many tokens we now have based on the time
		// passed since the previous request.
		refillsSince := uint(time.Since(b.time) / d)
		tokensAddedSince := refill * refillsSince
		currentTokens := b.tokens + tokensAddedSince

		// We don't have enough tokens. Return false.
		if currentTokens < 1 {
			//Con el false indicamos que se ha hecho throttle
			return false, "", nil
		}

		// If we've refilled our bucket, we can restart the clock.
		// Otherwise, we figure out when the most recent tokens were added.
		if currentTokens > max {
			b.time = time.Now()
			b.tokens = max - 1
		} else {
			deltaTokens := currentTokens - b.tokens
			deltaRefills := deltaTokens / refill
			deltaTime := time.Duration(deltaRefills) * d

			b.time = b.time.Add(deltaTime)
			b.tokens = currentTokens - 1
		}

		str, err := e(ctx)

		//Con el true indicamos que NO se ha hecho throttle
		return true, str, err
	}
}
