package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-srv/db"
	"rest-srv/models"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	stmt, err := db.Db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last inserted ID", http.StatusInternalServerError)
			return
		}
		addedTeachers[i].ID = int(lastID)
	}
	json.NewEncoder(w).Encode(addedTeachers)
	w.Header().Set("Content-Type", "application/json")
}

func TeacherHandler(w http.ResponseWriter, r *http.Request) {
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
