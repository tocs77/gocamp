package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"rest-srv/api/middlewares"
	"strconv"

	"golang.org/x/net/http2"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func teacherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "Handling POST teacher request...")
	case http.MethodGet:
		fmt.Fprintf(w, "Handling GET teacher request...")
	case http.MethodPut:
		fmt.Fprintf(w, "Handling PUT teacher request...")
	case http.MethodDelete:
		fmt.Fprintf(w, "Handling DELETE teacher request...")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func studentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "Handling POST student request...")
	case http.MethodGet:
		fmt.Fprintf(w, "Handling GET student request...")
	case http.MethodPut:
		fmt.Fprintf(w, "Handling PUT student request...")
	case http.MethodDelete:
		fmt.Fprintf(w, "Handling DELETE student request...")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "Handling POST exec request...")
	case http.MethodGet:
		fmt.Fprintf(w, "Handling GET exec request...")
	case http.MethodPut:
		fmt.Fprintf(w, "Handling PUT exec request...")
	case http.MethodDelete:
		fmt.Fprintf(w, "Handling DELETE exec request...")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	serverPort := 3000 // default port
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			serverPort = port
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/teachers", teacherHandler)
	mux.HandleFunc("/students", studentHandler)
	mux.HandleFunc("/execs", execHandler)

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
		Handler:   middlewares.ResponseTimMiddleware(middlewares.SecurityHeaders(middlewares.Cors(mux))),
	}

	http2.ConfigureServer(server, &http2.Server{})
	fmt.Println("Starting server on port ", serverPort)

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		fmt.Println("Error starting TLS server: ", err)
		os.Exit(1)
	}

}
