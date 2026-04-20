package handlers

import (
	pb "sch-grpc/proto/gen"
)

type Server struct {
	pb.UnimplementedTeachersServiceServer
	pb.UnimplementedStudentsServiceServer
	pb.UnimplementedExecsServiceServer
}
