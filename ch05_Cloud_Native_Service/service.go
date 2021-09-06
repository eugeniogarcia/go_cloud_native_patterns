package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var transact TransactionLogger

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Genera el evento para que se escriba en el log
	transact.WritePut(key, string(value))

	w.WriteHeader(http.StatusCreated)

	log.Printf("PUT key=%s value=%s\n", key, string(value))
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))

	log.Printf("GET key=%s\n", key)
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Genera el evento para que se escriba en el log
	transact.WriteDelete(key)

	log.Printf("DELETE key=%s\n", key)
}

//Rellena nuestro mapa con los datos encontrados en el transaction log
func initializeTransactionLog() error {
	var err error

	//Crea un logger en el archivo transactions.log
	//transact, err = NewFileTransactionLogger("transactions.log")

	//Crea un logger en postgress
	transact, err = NewPostgresTransactionLogger(PostgresDbParams{
		host:     "localhost",
		dbName:   "kvs",
		user:     "test",
		password: "hunter2",
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}

	//Lee el archivo y publica en estos canales lo que haya encontrado
	events, errors := transact.ReadEvents()
	count, ok, e := 0, true, Event{}

	//procesa los eventos, escribiendolos directamente en el mapa. Con esto lo que estamos haciendo es recuperar del log el estado del mapa, dejandolo como estaba en la última ejecución
	for ok && err == nil {
		select {
		case err, ok = <-errors:

		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = Delete(e.Key)
				count++
			case EventPut: // Got a PUT event!
				err = Put(e.Key, e.Value)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	//Ahora empezamos a actualizar el logger con los eventos que lleguen desde la API
	transact.Run()

	//Escribe en el log cualquier error que recibamos
	go func() {
		for err := range transact.Err() {
			log.Print(err)
		}
	}()

	return err
}

func main() {
	// Initializes the transaction log and loads existing data, if any.
	// Blocks until all data is read.
	err := initializeTransactionLog()
	if err != nil {
		panic(err)
	}

	// Create a new mux router
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	r.HandleFunc("/v1", notAllowedHandler)
	r.HandleFunc("/v1/{key}", notAllowedHandler)

	//Habilita https
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
	//Lo mismo, pero con http en lugar de https
	// log.Fatal(http.ListenAndServe(":8080", r))
}
