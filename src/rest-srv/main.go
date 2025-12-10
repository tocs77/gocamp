package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"rest-srv/api/middlewares"
	"rest-srv/api/router"
	"rest-srv/utility"
	"strconv"
	"time"

	"golang.org/x/net/http2"
)

func main() {

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
	rl := middlewares.NewRateLimiter(10, 2*time.Second)
	hpp := middlewares.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		WhiteList:                   []string{"name", "age", "address", "sortBy", "sortOrder"},
	}

	router := router.Router()
	middlewares := []utility.Middleware{
		middlewares.Hpp(hpp),
		middlewares.CompressionMiddleware,
		middlewares.SecurityHeaders,
		middlewares.ResponseTimMiddleware,
		rl.RateLimiterMiddleware,
		middlewares.Cors,
	}
	secureMux := utility.ApplyMiddlewares(router, middlewares...)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", serverPort),
		TLSConfig: tlsConfig,
		Handler:   secureMux,
	}

	http2.ConfigureServer(server, &http2.Server{})
	fmt.Println("Starting server on port ", serverPort)

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		fmt.Println("Error starting TLS server: ", err)
		os.Exit(1)
	}

}
