package utils

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

func ModelToBson[T any](model T, skipEmptyFields bool) (bson.M, error) {
	bsonFields := bson.M{}
	marshaledModel, err := bson.Marshal(model)
	if err != nil {
		return nil, HandleError(err, "failed to marshal model")
	}
	if err := bson.Unmarshal(marshaledModel, &bsonFields); err != nil {
		return nil, HandleError(err, "failed to unmarshal model")
	}
	delete(bsonFields, "_id")
	if skipEmptyFields {
		for key, value := range bsonFields {
			if IsZero(value) {
				delete(bsonFields, key)
			}
		}
	}
	return bsonFields, nil
}

func IsZero(value any) bool {
	return reflect.ValueOf(value).IsZero()
}
