package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"sch-grpc/cmd/grpcapi"
	mongodb "sch-grpc/internals/repositories"
	"sch-grpc/pkg/utils"
)

func main() {
	host := os.Getenv("MONGO_HOST")
	portStr := os.Getenv("MONGO_PORT")
	dbname := os.Getenv("MONGO_DBNAME")
	username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if host == "" || portStr == "" || dbname == "" {
		utils.HandleError(nil, "Failed to get MONGO_HOST, MONGO_PORT, or MONGO_DBNAME")
	}

	mongoPort, err := strconv.Atoi(portStr)
	if err != nil {
		utils.HandleError(err, "Failed to convert MONGO_PORT to int")
	}

	if err := mongodb.CreateMongoClient(context.Background(), host, mongoPort, dbname, username, password); err != nil {
		utils.HandleError(err, "Failed to create MongoDB client")
	}
	fmt.Println("Connected to MongoDB")
	defer mongodb.MongoClient.Disconnect(context.Background())

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		utils.HandleError(err, "Failed to convert SERVER_PORT to int")
	}
	grpcapi.RunServer(port)
}
