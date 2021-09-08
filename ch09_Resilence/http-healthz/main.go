package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/healthl", healthLivenessHandler)
	r.HandleFunc("/healths", healthShallowHandler)
	r.HandleFunc("/healthd", healthDeepHandler)
	r.HandleFunc("/healthc", healthCompoundHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
