package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

func getId(path string) (int, error) {
	p := strings.TrimPrefix(path, "/teachers")
	idStr := strings.Trim(p, "/")
	// Handle GET request for a specific teacher
	if idStr == "" {
		return 0, nil
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if id != 0 {
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

// teachers/{id}
func updateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "No ID", http.StatusBadRequest)
		return
	}

	// Verify teacher exists before updating
	var existingTeacher models.Teacher
	row := db.Db.QueryRow("SELECT * FROM teachers WHERE id = ?", id)
	err = row.Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error getting existing teacher", http.StatusInternalServerError)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Db.Prepare("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, id)
	if err != nil {
		http.Error(w, "Error updating data into database", http.StatusInternalServerError)
		return
	}

	// Set the ID from the URL path
	updatedTeacher.ID = id
	json.NewEncoder(w).Encode(updatedTeacher)
	w.Header().Set("Content-Type", "application/json")
}

func patchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "No ID", http.StatusBadRequest)
		return
	}

	var updatedFields map[string]any
	err = json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingTeacher models.Teacher
	row := db.Db.QueryRow("SELECT * FROM teachers WHERE id = ?", id)
	err = row.Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error getting existing teacher", http.StatusInternalServerError)
		return
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()
	for key, value := range updatedFields {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			jsonTag := field.Tag.Get("json")
			// Extract the field name from the JSON tag (remove ",omitempty" if present)
			jsonFieldName := strings.Split(jsonTag, ",")[0]
			if jsonFieldName == key && teacherVal.Field(i).CanSet() {
				teacherVal.Field(i).Set(reflect.ValueOf(value).Convert(teacherVal.Field(i).Type()))
				break
			}
		}
	}

	stmt, err := db.Db.Prepare("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, id)
	if err != nil {
		http.Error(w, "Error updating data into database", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existingTeacher)
	w.Header().Set("Content-Type", "application/json")
}

func TeacherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPut:
		updateTeacherHandler(w, r)
	case http.MethodPatch:
		patchTeacherHandler(w, r)
	case http.MethodDelete:
		fmt.Fprintf(w, "Handling DELETE teacher request...")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
