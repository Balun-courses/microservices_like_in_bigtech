package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	pb "github.com/moguchev/microservices_courcse/lecture_13/pkg/api/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	caCertFilePath = "./cert/ca-cert.crt"
)

func CreateClientTLSConfig(caCertFilePath string) (*tls.Config, error) {
	// Create a pool with the server certificate since it is not signed by a known CA
	caCert, err := os.ReadFile(caCertFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading server certificate: %s", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool, // aka: curl -v --cacert ./cert/server-cert.crt https://127.0.0.1:8443/hello
		InsecureSkipVerify: false,      // aka: curl -sL https://127.0.0.1:8443/hello --insecure
	}

	return tlsConfig, nil
}

func main() {
	tlsConfig, err := CreateClientTLSConfig(caCertFilePath)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("127.0.0.1:8082",
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewNotesServiceClient(conn)
	res, err := client.ListNotes(context.Background(), &pb.ListNotesRequest{})
	if err != nil {
		log.Fatal(err)
	}

	msg, _ := protojson.Marshal(res)
	log.Println("Msg:", string(msg))

}
