package mongodb

import (
	"context"
	"sch-grpc/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func deleteInDbByID(ctx context.Context, collectionName string, ids []string) ([]string, error) {
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, utils.HandleError(err, "failed to convert ID to object ID")
		}
		objectIds = append(objectIds, objectId)
	}
	result, err := MongoClient.Database("sch-db").Collection(collectionName).DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return nil, utils.HandleError(err, "failed to delete models in MongoDB")
	}
	if result.DeletedCount == 0 {
		return []string{}, nil
	}
	return ids, nil
}
