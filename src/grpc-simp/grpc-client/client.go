package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	pb "grpc-srv/protoc"
	farewellpb "grpc-srv/protoc/farewell"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	certFile := os.Getenv("CERT_FILE")
	if certFile == "" {
		certFile = "/workspace/cerificates/cert.pem"
	}
	serverName := os.Getenv("TLS_SERVER_NAME")
	if serverName == "" {
		serverName = "grpc-srv"
	}

	creds, err := credentials.NewClientTLSFromFile(certFile, serverName)
	if err != nil {
		log.Fatalf("Failed to load TLS cert: %v", err)
	}

	conn, err := grpc.NewClient("grpc-srv:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculateClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := client.Add(ctx, &pb.AddRequest{A: 10, B: 20})
	if err != nil {
		log.Fatalf("Failed to call Add: %v", err)
	}
	log.Printf("Add result: %d", res.GetSum())

	greeterClient := pb.NewGreeterClient(conn)
	greetRes, err := greeterClient.Greet(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Failed to call Greet: %v", err)
	}
	log.Printf("Greet result: %s", greetRes.GetMessage())

	// Generate Fibonacci

	fibonacciClient := pb.NewCalculateClient(conn)
	fibonacciRes, err := fibonacciClient.GenerateFibonacci(ctx, &pb.FibonacciRequest{N: 4})
	if err != nil {
		log.Fatalf("Failed to call GenerateFibonacci: %v", err)
	}
	for {
		fibonacci, err := fibonacciRes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive Fibonacci: %v", err)
		}
		log.Printf("Fibonacci result: %d", fibonacci.Number)
	}

	// Send Numbers
	numbersClient := pb.NewCalculateClient(conn)
	numbersStream, err := numbersClient.SendNumbers(ctx)
	if err != nil {
		log.Fatalf("Failed to call SendNumbers: %v", err)
	}
	defer numbersStream.CloseSend()

	for range 10 {
		err := numbersStream.Send(&pb.NumberRequest{Number: int32(rand.Intn(100))})
		if err != nil {
			log.Fatalf("Failed to send Number: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
		numbers, err := numbersStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive Numbers: %v", err)
		}
		log.Printf("Received number: %d, sum: %d", numbers.GetNumber(), numbers.GetSum())
	}

	// Chat
	chatClient := pb.NewCalculateClient(conn)
	chatStream, err := chatClient.Chat(ctx)
	if err != nil {
		log.Fatalf("Failed to call Chat: %v", err)
	}
	defer chatStream.CloseSend()

	go func() {
		for {
			message, err := chatStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive Chat: %v", err)
			}
			log.Printf("Received message: %s", message.GetMessage())
		}
	}()

	var messages = []string{"Hello, World!", "Hello, Go!", "Hello, Channels!"}

	for _, message := range messages {
		chatStream.Send(&pb.ChatMessage{Message: message})
		time.Sleep(1 * time.Second)
	}
	chatStream.CloseSend()

	// Farewell
	farewellClient := farewellpb.NewAufWiedersehenClient(conn)
	farewellRes, err := farewellClient.BigGoodBye(ctx, &farewellpb.GoodByeRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Failed to call BigGoodBye: %v", err)
	}
	log.Printf("BigGoodBye result: %s", farewellRes.GetMessage())
}
