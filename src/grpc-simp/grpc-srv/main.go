package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"

	pb "grpc-srv/protoc"
)

type server struct {
	pb.UnimplementedCalculateServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{Sum: req.A + req.B}, nil
}

func main() {
	port := 50051
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Printf("gRPC server listening on :%d", port)

	grpcServer := grpc.NewServer()
	pb.RegisterCalculateServer(grpcServer, &server{})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	fmt.Println("Server started on port ", port)
}
