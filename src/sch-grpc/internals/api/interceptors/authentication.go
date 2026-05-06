package interceptors

import (
	"context"
	"fmt"
	"sch-grpc/pkg/utils"
	"slices"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var skipMethods = []string{"/main.ExecsService/Login", "/main.ExecsService/ForgotPassword", "/main.ExecsService/ResetPassword"}

func AuthenticationInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	fmt.Println("AuthenticationInterceptor: ", info.FullMethod)
	if slices.Contains(skipMethods, info.FullMethod) {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header is not provided")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader[0], "Bearer "))
	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "token is not provided")
	}
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	newCtx := context.WithValue(ctx, "role", claims["role"].(string))
	newCtx = context.WithValue(newCtx, "exp", claims["exp"].(float64))
	newCtx = context.WithValue(newCtx, "username", claims["user"].(string))
	newCtx = context.WithValue(newCtx, "userId", claims["uid"].(string))
	return handler(newCtx, req)
}
