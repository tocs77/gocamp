package mongodb

import (
	"context"
	"errors"
	"sch-grpc/internals/models"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTeachers(ctx context.Context, teachers []*pb.Teacher) ([]*pb.Teacher, error) {
	newTeachers := make([]models.Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		var modelTeacher models.Teacher
		utils.MapStructFields(teacher, &modelTeacher)
		newTeachers = append(newTeachers, modelTeacher)
	}

	pbTeachers := make([]*pb.Teacher, 0, len(newTeachers))
	for _, teacher := range newTeachers {
		result, err := MongoClient.Database("sch-db").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, utils.HandleError(err, "failed to add teacher to MongoDB")
		}
		objId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, utils.HandleError(errors.New("failed to get object ID"), "failed to get object ID")
		}
		teacher.ID = objId.Hex()
		pbTeacher := &pb.Teacher{}
		utils.MapStructFields(teacher, pbTeacher)
		pbTeachers = append(pbTeachers, pbTeacher)
	}
	return pbTeachers, nil
}
