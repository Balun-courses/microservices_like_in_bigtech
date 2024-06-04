package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	caCertFilePath = "./cert/server-cert.crt"
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

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get("https://127.0.0.1:8443/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)
	msg, _ := io.ReadAll(resp.Body)
	log.Println("Msg:", string(msg))
}
