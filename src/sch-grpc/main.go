package main

import (
	"log"
	"os"
	"sch-grpc/cmd/grpcapi"
	"strconv"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Failed to convert SERVER_PORT to int: %v", err)
	}
	grpcapi.RunServer(port)
}
