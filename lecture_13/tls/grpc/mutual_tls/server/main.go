package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/moguchev/microservices_courcse/lecture_13/pkg/api/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	port           = ":8443"
	certFilePath   = "./cert/server-cert.crt"
	keyFilePath    = "./cert/server-key.key"
	caCertFilePath = "./cert/ca-cert.crt"
)

func CreateServerTLSConfig(caCertFilePath, certFile, keyFile string) (*tls.Config, error) {
	// Load server's certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509: %v", err)
	}

	// Load certificate of the CA who signed client's certificate
	clientCA, err := os.ReadFile(caCertFilePath)
	if err != nil {
		return nil, err
	}

	clientCAs := x509.NewCertPool()
	if !clientCAs.AppendCertsFromPEM(clientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	// RequireAndVerifyClientCert -
	// RequireAnyClientCert -
	// RequestClientCert -
	// NoClientCert -

	// Create tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // <--
		ClientCAs:    clientCAs,
	}

	return tlsConfig, nil
}

// server is used to implement pb.NotesServiceServer.
type server struct {
	pb.UnimplementedNotesServiceServer
}

func (s *server) ListNotes(_ context.Context, _ *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	return &pb.ListNotesResponse{}, nil
}

func main() {
	// create tls config
	tlsConfig, err := CreateServerTLSConfig(caCertFilePath, certFilePath, keyFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// create gRPC server
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	)

	// create router and register handlers
	srv := new(server)
	pb.RegisterNotesServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	// create listener
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// run server
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
