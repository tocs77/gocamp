package handlers

import (
	"fmt"
	"net/http"
)

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println(r.Form)
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
