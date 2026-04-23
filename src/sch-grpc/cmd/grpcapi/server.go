package grpcapi

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	"buf.build/go/protovalidate"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"sch-grpc/internals/api/handlers"
	"sch-grpc/pkg/utils"
	pb "sch-grpc/proto/gen"
)

var (
	globalCreds           credentials.TransportCredentials
	globalGatewayDialOpts []grpc.DialOption
	certFile              string
	keyFile               string
)

func init() {
	certFile = os.Getenv("CERT_FILE")
	if certFile == "" {
		fmt.Println("Error: CERT_FILE environment variable is required")
		os.Exit(1)
	}
	keyFile = os.Getenv("KEY_FILE")
	if keyFile == "" {
		fmt.Println("Error: KEY_FILE environment variable is required")
		os.Exit(1)
	}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}
	globalCreds = creds

	tlsServerName := os.Getenv("TLS_SERVER_NAME")
	if tlsServerName == "" {
		tlsServerName = "localhost"
	}

	caCertData, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Failed to read CERT_FILE: %v", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCertData) {
		log.Fatal("Failed to parse CERT_FILE as PEM CA certificate")
	}
	clientTLSConfig := &tls.Config{
		RootCAs:    certPool,
		ServerName: tlsServerName,
	}
	globalGatewayDialOpts = []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(clientTLSConfig)),
	}
}

// GatewayDialOptions returns dial options matching grpc-srv TLS trust setup (for gRPC-Gateway or other in-process clients).
func GatewayDialOptions() []grpc.DialOption {
	return globalGatewayDialOpts
}

func RunServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		utils.HandleError(err, "Failed to listen")
	}

	validator, err := protovalidate.New()
	if err != nil {
		utils.HandleError(err, "Failed to create validator")
	}

	interceptor := protovalidatemw.UnaryServerInterceptor(validator)

	grpcServer := grpc.NewServer(grpc.Creds(globalCreds), grpc.UnaryInterceptor(interceptor))
	pb.RegisterTeachersServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterStudentsServiceServer(grpcServer, &handlers.Server{})
	pb.RegisterExecsServiceServer(grpcServer, &handlers.Server{})
	reflection.Register(grpcServer)

	fmt.Printf("gRPC server listening on :%d (TLS)\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		utils.HandleError(err, "Failed to serve")
	}
}
