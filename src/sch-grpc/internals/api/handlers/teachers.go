package handlers

import (
	"context"

	mongodb "sch-grpc/internals/repositories/mongodb"
	pb "sch-grpc/proto/gen"

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
