package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"rest-srv/api/middlewares"
	"rest-srv/models"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

var teachers = make(map[int]models.Teacher)

var teachersMutex = &sync.Mutex{}
var nextTeacherID = 1

func init() {
	teachers[nextTeacherID] = models.Teacher{ID: nextTeacherID, FirstName: "John", LastName: "Doe", Class: "1A", Subject: "Math"}
	nextTeacherID++
	teachers[nextTeacherID] = models.Teacher{ID: nextTeacherID, FirstName: "Jane", LastName: "Smith", Class: "1B", Subject: "Science"}
	nextTeacherID++
	teachers[nextTeacherID] = models.Teacher{ID: nextTeacherID, FirstName: "Jim", LastName: "Beam", Class: "1C", Subject: "History"}
	nextTeacherID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers")
	idStr := strings.Trim(path, "/")
	// Handle GET request for a specific teacher
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		teacher, ok := teachers[id]
		if !ok {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(teacher)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	teachersList := make([]models.Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		if idStr != "" && strconv.Itoa(teacher.ID) != idStr {
			continue
		}
		teachersList = append(teachersList, teacher)
	}

	if sortBy != "" {
		sortByField := func(i, j int) bool {
			switch sortBy {
			case "id":
				return teachersList[i].ID < teachersList[j].ID

			case "first_name":
				return teachersList[i].FirstName < teachersList[j].FirstName
			case "last_name":
				return teachersList[i].LastName < teachersList[j].LastName
			case "class":
				return teachersList[i].Class < teachersList[j].Class
			case "subject":
				return teachersList[i].Subject < teachersList[j].Subject
			}
			return false
		}
		sort.SliceStable(teachersList, func(i, j int) bool {
			if sortOrder == "asc" {
				return sortByField(i, j)
			}
			return !sortByField(i, j)
		})
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{Status: "success", Count: len(teachersList), Data: teachersList}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	teachersMutex.Lock()
	defer teachersMutex.Unlock()
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range newTeachers {
		newTeachers[i].ID = nextTeacherID
		teachers[nextTeacherID] = newTeachers[i]
		nextTeacherID++
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{Status: "success", Count: len(newTeachers), Data: newTeachers}
	json.NewEncoder(w).Encode(response)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func teacherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodGet:
		getTeachersHandler(w, r)
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
	mux.HandleFunc("/teachers/", teacherHandler)
	mux.HandleFunc("/teachers", teacherHandler)
	mux.HandleFunc("/students", studentHandler)
	mux.HandleFunc("/execs", execHandler)
	mux.HandleFunc("/", rootHandler)

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

	middlewares := []Middleware{
		middlewares.Hpp(hpp),
		middlewares.CompressionMiddleware,
		middlewares.SecurityHeaders,
		middlewares.ResponseTimMiddleware,
		rl.RateLimiterMiddleware,
		middlewares.Cors,
	}
	secureMux := applyMiddlewares(mux, middlewares...)

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

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
