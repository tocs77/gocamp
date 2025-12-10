package handlers

import (
	"fmt"
	"net/http"
)

func ExecHandler(w http.ResponseWriter, r *http.Request) {
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
