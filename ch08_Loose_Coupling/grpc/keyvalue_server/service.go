package main

import (
	"context"
	"log"
	"net"

	pb "egsmartin.com/grpc/keyvalue"
	"google.golang.org/grpc"
)

//Creamos un tipo que representa la implementación de nuestro servidor. Para ello embebemos la definición del servidor que se genero al compilar el .proto. Esta implementación incluye definiciones para cada método "que no hacen nadas". En nuestro servidor crearemos una implementacion de estos metodos que "overloads" la implementación del tipo que hemos embebido

type server struct {
	pb.UnimplementedKeyValueServer
}

//Implementación de cada uno de los métodos del servidor
func (s *server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	value, err := Get(r.Key)

	return &pb.GetResponse{Value: value}, err
}

func (s *server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Received PUT key=%v value=%v", r.Key, r.Value)

	return &pb.PutResponse{}, Put(r.Key, r.Value)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Crea un servidor gRPC
	s := grpc.NewServer()

	//Asocia nuestro servicio al servidor
	pb.RegisterKeyValueServer(s, &server{})

	//Empieza a escuchar
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
