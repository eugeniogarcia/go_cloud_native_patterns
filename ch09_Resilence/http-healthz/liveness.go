package main

import "net/http"

//Test basico. Hay conectividad y el servicio esta vivo
func healthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
