package handlers

import (
	"context"

	"sch-grpc/internals/models"
	mongodb "sch-grpc/internals/repositories/mongodb"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//* AddTeachers

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

//* GetTeachers

func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.Teachers, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	// Filtering get filters from request
	filter, err := utils.BuildFilterForModel(models.Teacher{}, req.GetTeacher())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// Sorting get sort by from request
	sortOptions, err := utils.BuildSortForModel(models.Teacher{}, req.GetSortBy(), pb.Order_DESC)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	teachers, err := mongodb.GetTeachers(ctx, filter, sortOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &pb.Teachers{Teachers: teachers}, nil
}
