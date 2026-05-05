package mongodb

import (
	"context"
	"errors"
	"time"

	"sch-grpc/internals/models"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddExecs(ctx context.Context, execs []*pb.Exec) ([]*pb.Exec, error) {
	newExecs := make([]models.Exec, 0, len(execs))
	for _, exec := range execs {
		var modelExec models.Exec
		utils.MapStructFields(exec, &modelExec)
		hashedPassword, err := utils.HashPassword(exec.GetPassword())
		if err != nil {
			return nil, utils.HandleError(err, "failed to hash password")
		}
		modelExec.Password = hashedPassword
		modelExec.PasswordChangedAt = time.Now().Format(time.RFC3339)
		modelExec.UserCreatedAt = time.Now().Format(time.RFC3339)
		modelExec.InactiveStatus = false
		newExecs = append(newExecs, modelExec)
	}

	pbExecs := make([]*pb.Exec, 0, len(newExecs))
	for _, exec := range newExecs {
		result, err := MongoClient.Database("sch-db").Collection("execs").InsertOne(ctx, exec)
		if err != nil {
			return nil, utils.HandleError(err, "failed to add exec to MongoDB")
		}
		objID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, utils.HandleError(errors.New("failed to get object ID"), "failed to get object ID")
		}
		exec.ID = objID.Hex()

		pbExec := &pb.Exec{}
		utils.MapStructFields(exec, pbExec)
		pbExecs = append(pbExecs, pbExec)
	}

	return pbExecs, nil
}

func GetExecs(ctx context.Context, filter bson.M, sort bson.D) ([]*pb.ExecPublic, error) {
	collection := MongoClient.Database("sch-db").Collection("execs")
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, utils.HandleError(err, "failed to get execs from MongoDB")
	}
	defer cursor.Close(ctx)

	execs := make([]*pb.ExecPublic, 0)
	for cursor.Next(ctx) {
		var exec models.Exec
		if err := cursor.Decode(&exec); err != nil {
			return nil, utils.HandleError(err, "failed to decode exec from MongoDB")
		}

		pbExec := &pb.ExecPublic{}
		utils.MapStructFields(exec, pbExec)
		execs = append(execs, pbExec)
	}

	return execs, nil
}

func UpdateExecs(ctx context.Context, execs []*pb.Exec) ([]*pb.Exec, error) {
	updatedExecs, err := updateInDb[*pb.Exec, *pb.Exec](ctx, "execs", execs, bson.M{})
	if err != nil {
		return nil, utils.HandleError(err, "failed to update execs in MongoDB")
	}

	return updatedExecs, nil
}

func DeleteExecs(ctx context.Context, execIDs []string) ([]string, error) {
	deletedIDs, err := deleteInDbByID(ctx, "execs", execIDs)
	if err != nil {
		return nil, utils.HandleError(err, "failed to delete execs in MongoDB")
	}

	return deletedIDs, nil
}

func GetExecByUsername(ctx context.Context, username string) (*pb.Exec, error) {
	collection := MongoClient.Database("sch-db").Collection("execs")
	exec := models.Exec{}
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&exec)
	if err != nil {
		return nil, utils.HandleError(err, "failed to get exec by username from MongoDB")
	}
	pbExec := &pb.Exec{}
	utils.MapStructFields(exec, pbExec)
	return pbExec, nil
}

func DeactivateUsers(ctx context.Context, execIDs []string) ([]string, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range execIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, utils.HandleError(err, "failed to convert ID to object ID")
		}
		objectIDs = append(objectIDs, oid)
	}

	_, err := MongoClient.Database("sch-db").Collection("execs").UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}}, bson.M{"$set": bson.M{"inactive_status": true}})
	if err != nil {
		return nil, utils.HandleError(err, "failed to deactivate users in MongoDB")
	}

	return execIDs, nil
}

func GetExecByEmail(ctx context.Context, email string) (*pb.Exec, error) {
	collection := MongoClient.Database("sch-db").Collection("execs")
	exec := models.Exec{}
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&exec)
	if err != nil {
		return nil, utils.HandleError(err, "failed to get exec by email from MongoDB")
	}
	pbExec := &pb.Exec{}
	utils.MapStructFields(exec, pbExec)
	return pbExec, nil
}
