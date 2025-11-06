package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	serverPort := 3000 // default port
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			serverPort = port
		}
	}
	fmt.Println("Starting server on port ", serverPort)

	err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		os.Exit(1)
	}

}
