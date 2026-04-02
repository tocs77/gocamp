package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "grpc-srv/protoc"
	farewellpb "grpc-srv/protoc/farewell"
)

type server struct {
	pb.UnimplementedCalculateServer
	pb.UnimplementedGreeterServer
	farewellpb.UnimplementedAufWiedersehenServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{Sum: req.A + req.B}, nil
}

func (s *server) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func (s *server) BigGoodBye(ctx context.Context, req *farewellpb.GoodByeRequest) (*farewellpb.GoodByeResponse, error) {
	return &farewellpb.GoodByeResponse{Message: "Goodbye, " + req.Name}, nil
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

	// SSL certificate and key (required)
	certFile := os.Getenv("CERT_FILE")
	if certFile == "" {
		fmt.Println("Error: CERT_FILE environment variable is required")
		os.Exit(1)
	}
	keyFile := os.Getenv("KEY_FILE")
	if keyFile == "" {
		fmt.Println("Error: KEY_FILE environment variable is required")
		os.Exit(1)
	}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	s := &server{}
	pb.RegisterCalculateServer(grpcServer, s)
	pb.RegisterGreeterServer(grpcServer, s)
	farewellpb.RegisterAufWiedersehenServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	fmt.Println("Server started on port ", port)
}
