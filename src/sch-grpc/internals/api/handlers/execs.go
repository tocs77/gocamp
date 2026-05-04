package handlers

import (
	"context"
	"fmt"
	"time"

	"sch-grpc/internals/models"
	mongodb "sch-grpc/internals/repositories/mongodb"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

//* Login

func (s *Server) Login(ctx context.Context, req *pb.ExecLoginRequest) (*pb.ExecLoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	exec, err := mongodb.GetExecByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if exec.InactiveStatus {
		return nil, status.Error(codes.Unauthenticated, "user is inactive")
	}
	valid, err := utils.ComparePassword(exec.GetPassword(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}
	token, err := utils.SignToken(exec.GetId(), exec.GetUsername(), exec.GetRole())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.ExecLoginResponse{Status: true, Token: token}, nil
}

//* UpdatePassword

func (s *Server) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	fmt.Println("UpdatePassword request received")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is required")
	}
	if req.GetCurrentPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. Current password field is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. ID field is not a valid ObjectID")
	}

	exec := models.Exec{}
	err = mongodb.MongoClient.Database("sch-db").Collection("execs").FindOne(ctx, bson.M{"_id": objectID}).Decode(&exec)
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error finding exec").Error())
	}
	if exec.InactiveStatus {
		return nil, status.Error(codes.Unauthenticated, "user is inactive")
	}
	valid, err := utils.ComparePassword(exec.Password, req.GetCurrentPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error comparing passwords").Error())
	}
	if !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid current password")
	}
	hashedPassword, err := utils.HashPassword(req.GetNewPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error hashing password").Error())
	}

	_, err = mongodb.MongoClient.Database("sch-db").Collection("execs").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"password": hashedPassword, "password_changed_at": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error updating exec").Error())
	}
	token, err := utils.SignToken(exec.ID, exec.Username, exec.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error signing token").Error())
	}
	return &pb.UpdatePasswordResponse{PasswordUpdated: true, Token: token}, nil
}

//* Deactivate user

func (s *Server) DeactivateUser(ctx context.Context, req *pb.ExecsIds) (*pb.Confirmation, error) {
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

	_, err := mongodb.DeactivateUsers(ctx, execIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error deactivating users").Error())
	}

	return &pb.Confirmation{Confirmation: true}, nil
}
