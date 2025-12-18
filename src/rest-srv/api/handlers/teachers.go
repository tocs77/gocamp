package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"rest-srv/db"
	"rest-srv/models"
	"strconv"
	"strings"
)

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
		row := db.Db.QueryRow("SELECT * FROM teachers WHERE id = ?", id)
		if row == nil {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		var teacher models.Teacher
		err = row.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Invalid record in database", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(teacher)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")

	// Validate and sanitize sortBy to prevent SQL injection
	validColumns := map[string]bool{
		"id":         true,
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	var query string
	if validColumns[sortBy] {
		// Validate sortOrder (only allow ASC or DESC)
		if sortOrder == "" {
			sortOrder = "ASC"
		}
		sortOrder = strings.ToUpper(sortOrder)
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "ASC"
		}

		// Build query with validated column name and sort order
		query = fmt.Sprintf("SELECT * FROM teachers ORDER BY %s %s", sortBy, sortOrder)
	} else {
		// Invalid or empty sortBy, query without ORDER BY
		query = "SELECT * FROM teachers"
	}

	teachersList := make([]models.Teacher, 0)
	rows, err := db.Db.Query(query)
	if err != nil {
		fmt.Println("Error querying teachers: ", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			http.Error(w, "Invalid record in database", http.StatusInternalServerError)
			return
		}
		teachersList = append(teachersList, teacher)
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
