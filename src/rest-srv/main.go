package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Person struct {
	Name string
	Age  int
}

var people = map[string]Person{
	"1": {Name: "John", Age: 30},
	"2": {Name: "Jane", Age: 25},
	"3": {Name: "Jim", Age: 35},
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/person", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		person, ok := people[id]
		if !ok {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(person)
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
