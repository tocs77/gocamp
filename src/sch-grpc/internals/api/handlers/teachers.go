package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"sch-grpc/internals/models"
	mongodb "sch-grpc/internals/repositories/mongodb"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	for _, teacher := range req.GetTeachers() {
		if teacher.GetId() != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is not allowed")
		}
	}

	pbTeachers, err := mongodb.AddTeachers(ctx, req.GetTeachers())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &pb.Teachers{Teachers: pbTeachers}, nil
}

func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.Teachers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	// Filtering get filters from request
	filter, err := buildFilterForTeacher(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fmt.Println("filter", filter)
	// Sorting get sort by from request
	sortOptions, err := buildSortForTeacher(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fmt.Println("sortOptions", sortOptions)
	return nil, nil
}

func buildFilterForTeacher(req *pb.GetTeachersRequest) (bson.M, error) {
	return utils.BuildFilterForModel(models.Teacher{}, req.GetTeacher())
}

func buildSortForTeacher(req *pb.GetTeachersRequest) (bson.D, error) {
	sortOptions := bson.D{}
	modelType := reflect.TypeFor[models.Teacher]()
	allowedSortFields := make(map[string]struct{}, modelType.NumField())

	for i := 0; i < modelType.NumField(); i++ {
		bsonTag := modelType.Field(i).Tag.Get("bson")
		columnName := strings.Split(bsonTag, ",")[0]
		if columnName != "" {
			allowedSortFields[columnName] = struct{}{}
		}
	}

	for _, sortBy := range req.GetSortBy() {
		if _, ok := allowedSortFields[sortBy.GetField()]; !ok {
			continue
		}

		order := 1
		if sortBy.GetOrder() == pb.Order_DESC {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{Key: sortBy.GetField(), Value: order})
	}
	return sortOptions, nil
}
