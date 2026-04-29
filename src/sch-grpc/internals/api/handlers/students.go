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

//* AddStudents

func (s *Server) AddStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	for _, student := range req.GetStudents() {
		if student.GetId() != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is not allowed")
		}
	}

	pbStudents, err := mongodb.AddStudents(ctx, req.GetStudents())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.Students{Students: pbStudents}, nil
}

//* GetStudents

func (s *Server) GetStudents(ctx context.Context, req *pb.GetStudentsRequest) (*pb.Students, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	filter, err := utils.BuildFilterForModel(models.Student{}, req.GetStudents())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sortOptions, err := utils.BuildSortForModel(models.Student{}, req.GetSortBy(), pb.Order_DESC)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	students, err := mongodb.GetStudents(ctx, filter, sortOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.Students{Students: students}, nil
}

//* UpdateStudents

func (s *Server) UpdateStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	for _, student := range req.GetStudents() {
		if student.GetId() == "" {
			return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is required")
		}
	}

	updatedStudents, err := mongodb.UpdateStudents(ctx, req.GetStudents())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.Students{Students: updatedStudents}, nil
}

//* DeleteStudents

func (s *Server) DeleteStudents(ctx context.Context, req *pb.StudentsIds) (*pb.DeleteStudentsConfirmation, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	for _, student := range req.GetIds() {
		if student.GetId() == "" {
			return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is required")
		}
	}

	studentIDs := make([]string, 0, len(req.GetIds()))
	for _, student := range req.GetIds() {
		studentIDs = append(studentIDs, student.GetId())
	}

	deletedIDs, err := mongodb.DeleteStudents(ctx, studentIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	deletedStudentIDs := make([]*pb.StudentId, 0, len(deletedIDs))
	for _, deletedID := range deletedIDs {
		deletedStudentIDs = append(deletedStudentIDs, &pb.StudentId{Id: deletedID})
	}

	return &pb.DeleteStudentsConfirmation{
		Status:     "success",
		DeletedIds: deletedStudentIDs,
	}, nil
}
