package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/net/http2"
)

func main() {

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		logRequestDetails(r)
		logTlsVersion(r.TLS.Version)
		fmt.Fprintf(w, "Handling orders...")
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		logRequestDetails(r)
		logTlsVersion(r.TLS.Version)
		fmt.Fprintf(w, "Handling users...")
	})

	serverPort := 3000 // default port
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			serverPort = port
		}
	}

	//Load the SSL certificate and key

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading SSL certificate and key: ", err)
		os.Exit(1)
	}
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", serverPort),
		TLSConfig: tlsConfig,
	}

	http2.ConfigureServer(server, &http2.Server{})
	fmt.Println("Starting server on port ", serverPort)

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		fmt.Println("Error starting TLS server: ", err)
		os.Exit(1)
	}

}

func logRequestDetails(r *http.Request) {
	httpVersion := r.Proto
	fmt.Printf("HTTP Version: %s\n", httpVersion)
}

func logTlsVersion(version uint16) {
	switch version {
	case tls.VersionTLS10:
		fmt.Printf("TLS Version: TLS 1.0\n")
	case tls.VersionTLS11:
		fmt.Printf("TLS Version: TLS 1.1\n")
	case tls.VersionTLS12:
		fmt.Printf("TLS Version: TLS 1.2\n")
	case tls.VersionTLS13:
		fmt.Printf("TLS Version: TLS 1.3\n")
	default:
		fmt.Printf("TLS Version: Unknown\n")
	}

}
