package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-srv/db"
	"rest-srv/models"
	"strconv"
	"strings"
)

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	student, err := db.GetStudentById(id)
	if err != nil {
		if err.Error() == "student not found" {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to retrieve student", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
	w.Header().Set("Content-Type", "application/json")
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	sortByParams := r.URL.Query()["sortBy"]
	queryParams := r.URL.Query()

	// Parse filter parameters (format: field=value)
	filters := make(map[string]string)
	for field, values := range queryParams {
		// Skip sortBy parameter as it's handled separately
		if field == "sortBy" {
			continue
		}

		// Use the first value if multiple are provided
		if len(values) > 0 && values[0] != "" {
			filters[field] = values[0]
		}
	}

	studentsList, err := db.GetStudents(filters, sortByParams)
	if err != nil {
		http.Error(w, "unable to retrieve students", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{Status: "success", Count: len(studentsList), Data: studentsList}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	var newStudents []models.Student
	err := json.NewDecoder(r.Body).Decode(&newStudents)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for _, student := range newStudents {
		err = student.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedStudents, err := db.AddStudents(newStudents)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(addedStudents)
	w.Header().Set("Content-Type", "application/json")
}

// students/{id}
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err = updatedStudent.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedStudent, err = db.UpdateStudent(id, updatedStudent)
	if err != nil {
		if err.Error() == "student not found" {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update student", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedStudent)
	w.Header().Set("Content-Type", "application/json")
}

func PatchStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var updatedFields map[string]any
	err = json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedStudent, err := db.PatchStudent(id, updatedFields)
	if err != nil {
		if err.Error() == "student not found" {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update student", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedStudent)
	w.Header().Set("Content-Type", "application/json")
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedStudents, err := db.PatchStudents(updates)
	if err != nil {
		if err.Error() == "student not found" || strings.Contains(err.Error(), "student not found") {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update students", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedStudents)
	w.Header().Set("Content-Type", "application/json")
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	deletedStudent, err := db.DeleteStudent(id)
	if err != nil {
		if err.Error() == "student not found" {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to delete student", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		ID      int    `json:"id"`
	}{Status: "success", Message: "Student deleted successfully", ID: deletedStudent.ID}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deletedStudents, err := db.DeleteStudents(ids)
	if err != nil {
		if strings.Contains(err.Error(), "student not found") {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to delete students", http.StatusInternalServerError)
		return
	}

	deletedIDs := make([]int, len(deletedStudents))
	for i, student := range deletedStudents {
		deletedIDs[i] = student.ID
	}

	json.NewEncoder(w).Encode(struct {
		Status     string `json:"status"`
		Message    string `json:"message"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{Status: "success", Message: "Students deleted successfully", DeletedIDs: deletedIDs})
	w.Header().Set("Content-Type", "application/json")
}
