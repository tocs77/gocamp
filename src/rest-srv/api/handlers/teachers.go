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

func checkValidField(field string) bool {
	// Validate and sanitize sortBy to prevent SQL injection
	validColumns := map[string]bool{
		"id":         true,
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validColumns[field]
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

	sortByParams := r.URL.Query()["sortBy"]
	queryParams := r.URL.Query()

	var query string
	var orderByClauses []string
	var whereClauses []string
	var filterValues []interface{}

	// Parse filter parameters (format: field=value)
	for field, values := range queryParams {
		// Skip sortBy parameter as it's handled separately
		if field == "sortBy" {
			continue
		}

		// Validate field
		if !checkValidField(field) {
			continue // Skip invalid field
		}

		// Use the first value if multiple are provided
		if len(values) > 0 && values[0] != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			filterValues = append(filterValues, values[0])
		}
	}

	// Parse each sortBy parameter (format: field:asc or field:desc)
	for _, sortParam := range sortByParams {
		parts := strings.Split(sortParam, ":")
		if len(parts) != 2 {
			continue // Skip invalid format
		}

		field := strings.TrimSpace(parts[0])
		order := strings.TrimSpace(strings.ToUpper(parts[1]))

		// Validate field
		if !checkValidField(field) {
			continue // Skip invalid field
		}

		// Validate order (only allow ASC or DESC)
		if order != "ASC" && order != "DESC" {
			order = "ASC" // Default to ASC if invalid
		}

		// Add to order by clauses
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", field, order))
	}

	// Build query with WHERE and ORDER BY clauses
	query = "SELECT * FROM teachers"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	if len(orderByClauses) > 0 {
		query += " ORDER BY " + strings.Join(orderByClauses, ", ")
	}

	teachersList := make([]models.Teacher, 0)
	var rows *sql.Rows
	var err error

	// Execute query with parameterized values
	if len(filterValues) > 0 {
		rows, err = db.Db.Query(query, filterValues...)
	} else {
		rows, err = db.Db.Query(query)
	}
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
