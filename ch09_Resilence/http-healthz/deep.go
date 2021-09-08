package main

import (
	"context"
	"net/http"
	"time"
)

var service Service

//Comprueba que se pueda llamar al servicio
func healthDeepHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the context from the request and add a 5-second timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Can the service execute a key query against the database?
	if err := service.GetUser(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Service struct{}

func (s Service) GetUser(ctx context.Context) error {
	// An imaginary function that executes a simple database query.
	if err := HealthCheck(ctx); err != nil {
		return err
	}

	return nil
}

func HealthCheck(ctx context.Context) error {
	time.Sleep(500 * time.Millisecond)
	return ctx.Err()
}
