package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func CreateMongoClient(ctx context.Context, host string, port int, dbname string, username string, password string) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d/%s", host, port, dbname))
	if username != "" && password != "" {
		clientOptions.SetAuth(options.Credential{
			Username:   username,
			Password:   password,
			AuthSource: "admin",
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	MongoClient = client
	return nil
}
