package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

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

func (s *server) GenerateFibonacci(req *pb.FibonacciRequest, stream pb.Calculate_GenerateFibonacciServer) error {
	n := int(req.N)
	a, b := 0, 1
	for range n {
		err := stream.Send(&pb.FibonacciResponse{Number: int32(a)})
		if err != nil {
			return err
		}
		a, b = b, a+b
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *server) SendNumbers(stream pb.Calculate_SendNumbersServer) error {
	sum := 0
	for {
		number, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sum += int(number.Number)
		err = stream.Send(&pb.NumberResponse{Number: number.Number, Sum: int32(sum)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func (s *server) BigGoodBye(ctx context.Context, req *farewellpb.GoodByeRequest) (*farewellpb.GoodByeResponse, error) {
	return &farewellpb.GoodByeResponse{Message: "Goodbye, " + req.Name}, nil
}

func (s *server) Chat(stream pb.Calculate_ChatServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			return err
		}
		log.Printf("Received message: %s", message.Message)
		stream.Send(&pb.ChatMessage{Message: "Received message: " + message.Message})
	}
	return nil
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
