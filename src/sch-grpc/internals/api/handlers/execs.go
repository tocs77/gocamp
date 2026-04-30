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

//* AddExecs

func (s *Server) AddExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	for _, exec := range req.GetExecs() {
		if exec.GetId() != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is not allowed")
		}
	}

	pbExecs, err := mongodb.AddExecs(ctx, req.GetExecs())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.Execs{Execs: pbExecs}, nil
}

//* GetExecs

func (s *Server) GetExecs(ctx context.Context, req *pb.GetExecsRequest) (*pb.ExecsPublic, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	filter, err := utils.BuildFilterForModel(models.ExecPublic{}, req.GetExecs())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sortOptions, err := utils.BuildSortForModel(models.ExecPublic{}, req.GetSortBy(), pb.Order_DESC)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	execs, err := mongodb.GetExecs(ctx, filter, sortOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.ExecsPublic{Execs: execs}, nil
}

//* UpdateExecs

func (s *Server) UpdateExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	for _, exec := range req.GetExecs() {
		if exec.GetId() == "" {
			return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is required")
		}
	}

	updatedExecs, err := mongodb.UpdateExecs(ctx, req.GetExecs())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.Execs{Execs: updatedExecs}, nil
}

//* DeleteExecs

func (s *Server) DeleteExecs(ctx context.Context, req *pb.ExecsIds) (*pb.DeleteExecsConfirmation, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	for _, exec := range req.GetIds() {
		if exec.GetId() == "" {
			return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is required")
		}
	}

	execIDs := make([]string, 0, len(req.GetIds()))
	for _, exec := range req.GetIds() {
		execIDs = append(execIDs, exec.GetId())
	}

	deletedIDs, err := mongodb.DeleteExecs(ctx, execIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	deletedExecIDs := make([]*pb.ExecId, 0, len(deletedIDs))
	for _, deletedID := range deletedIDs {
		deletedExecIDs = append(deletedExecIDs, &pb.ExecId{Id: deletedID})
	}

	return &pb.DeleteExecsConfirmation{
		Status:     "success",
		DeletedIds: deletedExecIDs,
	}, nil
}
