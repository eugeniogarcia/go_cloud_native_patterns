package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	pb "egsmartin.com/grpc/keyvalue"
	"google.golang.org/grpc"
)

func main() {
	// Crea un contexto con time-out
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Opciones de conexión con el servidor gRPC
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	//Establece la conexión con el servidor
	conn, err := grpc.DialContext(ctx, "localhost:50051", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Crea una instancia del cliente, usando la conexión que acabamos de crear
	client := pb.NewKeyValueClient(conn)

	var action, key, value string

	//Llama a los distintos métodos del cliente, usando el contexto que tiene definido un timeout de 1 seg
	// Expect something like "set foo bar"
	if len(os.Args) > 2 {
		action, key = os.Args[1], os.Args[2]
		value = strings.Join(os.Args[3:], " ")
	}

	// Call client.Get() or client.Put() as appropriate.
	switch action {
	case "get":
		r, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			log.Fatalf("could not get value for key %s: %v\n", key, err)
		}
		log.Printf("Get %s returns: %s", key, r.Value)

	case "put":
		_, err := client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("could not get put key %s: %v\n", key, err)
		}
		log.Printf("Put %s", key)

	default:
		log.Fatalf("Syntax: go run [get|put] KEY VALUE...")
	}
}
