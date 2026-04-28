package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sch-grpc/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func updateInDb[T any, Result any](ctx context.Context, collectionName string, model []T, filter bson.M) ([]Result, error) {
	collection := MongoClient.Database("sch-db").Collection(collectionName)
	succesFullyUpdatedModels := make([]Result, 0, len(model))

	for _, item := range model {
		modelID, err := extractModelID(item)
		if err != nil {
			return nil, utils.HandleError(err, "failed to get model ID")
		}
		objectId, err := primitive.ObjectIDFromHex(modelID)
		if err != nil {
			return nil, utils.HandleError(err, "failed to convert object ID")
		}
		updateFields, err := utils.ModelToBson(item, true)
		if err != nil {
			return nil, utils.HandleError(err, "failed to convert model to bson")
		}
		var updatedModel Result
		err = collection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": objectId},
			bson.M{"$set": updateFields},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&updatedModel)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, utils.HandleError(errors.New("model not found"), "model not found")
			}
			return nil, utils.HandleError(err, "failed to update model in MongoDB")
		}
		succesFullyUpdatedModels = append(succesFullyUpdatedModels, updatedModel)
	}
	return succesFullyUpdatedModels, nil
}

func extractModelID(model any) (string, error) {
	modelValue := reflect.ValueOf(model)
	if !modelValue.IsValid() {
		return "", errors.New("model is invalid")
	}
	if modelValue.Kind() == reflect.Pointer {
		if modelValue.IsNil() {
			return "", errors.New("model is nil")
		}
		modelValue = modelValue.Elem()
	}
	if modelValue.Kind() != reflect.Struct {
		return "", fmt.Errorf("model kind must be struct, got %s", modelValue.Kind())
	}

	field := modelValue.FieldByName("ID")
	if !field.IsValid() {
		field = modelValue.FieldByName("Id")
	}
	if !field.IsValid() || field.Kind() != reflect.String || field.Len() == 0 {
		return "", errors.New("model ID field is missing or empty")
	}
	return field.String(), nil
}
