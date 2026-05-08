package handlers

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"sch-grpc/internals/models"
	mongodb "sch-grpc/internals/repositories/mongodb"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"

	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"
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
	utils.JWTStorage.AddToken(token)
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

//* Forgot Password

func (s *Server) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. Email field is required")
	}

	exec, err := mongodb.GetExecByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error retrieving exec").Error())
	}
	if exec.InactiveStatus {
		return nil, status.Error(codes.Unauthenticated, "user is inactive")
	}
	tokenbytes := make([]byte, 32)
	rand.Read(tokenbytes)
	token := hex.EncodeToString(tokenbytes)
	hashedToken := utils.PasswordResetTokenHashFromRaw(tokenbytes)
	expires := time.Now().Add(time.Hour * 24).Format(time.RFC3339)
	if err := mongodb.SetExecPasswordResetFields(ctx, exec.GetId(), hashedToken, expires); err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error updating exec").Error())
	}
	resetPasswordURL := fmt.Sprintf("https://localhost:%s/execs/reset-password/reset/%s", os.Getenv("EXPOSE_PORT"), token)
	return &pb.ForgotPasswordResponse{Confirmation: true, Message: "Forgot password?. Reset your password using following link: " + resetPasswordURL}, nil
}

//* Reset Password

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.Confirmation, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.GetResetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. Reset code field is required")
	}
	if req.GetNewPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. New password field is required")
	}
	if req.GetConfirmPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. Confirm password field is required")
	}
	if req.GetNewPassword() != req.GetConfirmPassword() {
		return nil, status.Error(codes.InvalidArgument, "request is in invalid format. New password and confirm password do not match")
	}
	exec, err := mongodb.GetExecByPasswordResetToken(ctx, req.GetResetCode())
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, status.Error(codes.Unauthenticated, "invalid or expired reset link")
		case errors.Is(err, utils.ErrInvalidPasswordResetToken):
			return nil, status.Error(codes.InvalidArgument, "invalid reset code")
		default:
			return nil, status.Error(codes.Internal, utils.HandleError(err, "error retrieving exec").Error())
		}
	}
	if exec.InactiveStatus {
		return nil, status.Error(codes.Unauthenticated, "user is inactive")
	}
	expiresAt, err := time.Parse(time.RFC3339, exec.GetPasswordTokenExpires())
	if err != nil || time.Now().After(expiresAt) {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired reset link")
	}
	hashedPassword, err := utils.HashPassword(req.GetNewPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error hashing password").Error())
	}
	objectID, err := primitive.ObjectIDFromHex(exec.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}
	_, err = mongodb.MongoClient.Database("sch-db").Collection("execs").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{
		"password":               hashedPassword,
		"password_changed_at":    time.Now().Format(time.RFC3339),
		"password_reset_token":   "",
		"password_token_expires": "",
	}})
	if err != nil {
		return nil, status.Error(codes.Internal, utils.HandleError(err, "error updating password").Error())
	}
	return &pb.Confirmation{Confirmation: true}, nil
}

//* Logout

func (s *Server) Logout(ctx context.Context, req *pb.EmptyRequest) (*pb.ExecLogoutResponse, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}
	authHeader := metadata.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header is not provided")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader[0], "Bearer "))
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is not provided")
	}
	utils.JWTStorage.DeleteToken(token)
	return &pb.ExecLogoutResponse{LoggedOut: true}, nil
}
