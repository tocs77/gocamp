package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, err := grpc.NewClient("grpc-srv:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	log.Println("Client created")
	defer client.Close()

}
