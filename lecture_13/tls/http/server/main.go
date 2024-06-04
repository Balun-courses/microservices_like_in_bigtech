package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	// "golang.org/x/net/http2"
)

const (
	port         = ":8443"
	certFilePath = "./cert/server-cert.crt"
	keyFilePath  = "./cert/server-key.key"
)

func CreateServerTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	// Load server's certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509: %v", err)
	}

	// Create tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}

	return tlsConfig, nil
}

func main() {
	// create tls config
	tlsConfig, err := CreateServerTLSConfig(certFilePath, keyFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// create router and register handlers
	router := http.NewServeMux()
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, TLS!"))
	})

	// create HTTP server
	server := &http.Server{
		Addr:      port,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	// 	create listener and run server
	log.Printf("Listening on [::]%s ...", port)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// Unsafe (no TLS):
	// curl -sL https://127.0.0.1:8443/hello --insecure

	// Safe (With TLS):
	// curl -v --cacert ./cert/server-cert.crt https://127.0.0.1:8443/hello
	// curl -v --cacert ./cert/ca-cert.crt https://127.0.0.1:8443/hello
}
