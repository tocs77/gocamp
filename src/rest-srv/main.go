package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"rest-srv/api/middlewares"
	"rest-srv/api/router"
	"rest-srv/db"
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

	// Database connection parameters from environment - all required
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		fmt.Println("Error: DB_HOST environment variable is required")
		os.Exit(1)
	}

	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		fmt.Println("Error: DB_PORT environment variable is required")
		os.Exit(1)
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		fmt.Printf("Error: DB_PORT must be a valid integer: %v\n", err)
		os.Exit(1)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		fmt.Println("Error: DB_USER environment variable is required")
		os.Exit(1)
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		fmt.Println("Error: DB_PASSWORD environment variable is required")
		os.Exit(1)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		fmt.Println("Error: DB_NAME environment variable is required")
		os.Exit(1)
	}

	// Connect to database
	if err := db.ConnectDb(dbUser, dbPassword, dbHost, dbPort, dbName); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Db.Close()

	// Test database connection
	if err := db.Db.Ping(); err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Database connection established successfully")

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
		WhiteList:                   []string{"name", "age", "address", "sortBy", "sortOrder", "id", "first_name", "last_name", "email", "class", "subject"},
	}

	router := router.MainRouter()
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
