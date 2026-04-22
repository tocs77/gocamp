package handlers

import (
	"context"
	"errors"

	"sch-grpc/internals/models"
	mongodb "sch-grpc/internals/repositories"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	newTeachers := make([]models.Teacher, 0, len(req.GetTeachers()))
	for _, teacher := range req.GetTeachers() {
		var modelTeacher models.Teacher
		utils.MapPBModelToModel(teacher, &modelTeacher)
		newTeachers = append(newTeachers, modelTeacher)
	}

	pbTeachers := make([]*pb.Teacher, 0, len(newTeachers))
	for _, teacher := range newTeachers {
		result, err := mongodb.MongoClient.Database("sch-db").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, utils.HandleError(err, "failed to add teacher to MongoDB")
		}
		objId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, errors.New("failed to get object ID")
		}
		teacher.ID = objId.Hex()
		pbTeacher := &pb.Teacher{}
		utils.MapModelToPB(teacher, pbTeacher)
		pbTeachers = append(pbTeachers, pbTeacher)
	}
	return &pb.Teachers{Teachers: pbTeachers}, nil
}
