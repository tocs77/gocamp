package handlers

import (
	"encoding/json"
	"net/http"
	"rest-srv/db"
	"rest-srv/models"
	"rest-srv/utility"
	"strconv"
	"strings"
)

func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	teacher, err := db.GetTeacherById(id)
	if err != nil {
		if err.Error() == "teacher not found" {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to retrieve teacher", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teacher)
	w.Header().Set("Content-Type", "application/json")

}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
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

	teachersList, err := db.GetTeachers(filters, sortByParams)
	if err != nil {
		http.Error(w, "unable to retrieve teachers", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{Status: "success", Count: len(teachersList), Data: teachersList}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for _, teacher := range newTeachers {
		err = teacher.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedTeachers, err := db.AddTeachers(newTeachers)
	if err != nil {
		http.Error(w, "unable to add teachers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(addedTeachers)
	w.Header().Set("Content-Type", "application/json")
}

// teachers/{id}
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err = updatedTeacher.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedTeacher, err = db.UpdateTeacher(id, updatedTeacher)
	if err != nil {
		if err.Error() == "teacher not found" {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update teacher", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTeacher)
	w.Header().Set("Content-Type", "application/json")
}

func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	updatedTeacher, err := db.PatchTeacher(id, updatedFields)
	if err != nil {
		if err.Error() == "teacher not found" {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update teacher", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTeacher)
	w.Header().Set("Content-Type", "application/json")
}

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedTeachers, err := db.PatchTeachers(updates)
	if err != nil {
		if err.Error() == "teacher not found" || strings.Contains(err.Error(), "teacher not found") {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update teachers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTeachers)
	w.Header().Set("Content-Type", "application/json")
}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	deletedTeacher, err := db.DeleteTeacher(id)
	if err != nil {
		if err.Error() == "teacher not found" {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to delete teacher", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		ID      int    `json:"id"`
	}{Status: "success", Message: "Teacher deleted successfully", ID: deletedTeacher.ID}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deletedTeachers, err := db.DeleteTeachers(ids)
	if err != nil {
		if strings.Contains(err.Error(), "teacher not found") {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to delete teachers", http.StatusInternalServerError)
		return
	}

	deletedIDs := make([]int, len(deletedTeachers))
	for i, teacher := range deletedTeachers {
		deletedIDs[i] = teacher.ID
	}

	json.NewEncoder(w).Encode(struct {
		Status     string `json:"status"`
		Message    string `json:"message"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{Status: "success", Message: "Teachers deleted successfully", DeletedIDs: deletedIDs})
	w.Header().Set("Content-Type", "application/json")
}

func GetTeacherStudentsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	students, err := db.GetTeacherStudents(id)
	if err != nil {
		http.Error(w, "unable to retrieve students", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{Status: "success", Count: len(students), Data: students})
	w.Header().Set("Content-Type", "application/json")
}

func GetTeacherStudentsCountHandler(w http.ResponseWriter, r *http.Request) {
	allowedRoles := []string{"admin", "exec", "manager"}
	authorized, err := utility.AuthorizeUser(r.Context().Value(utility.ContextKey("role")).(string), allowedRoles...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !authorized {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	count, err := db.GetTeacherStudentsCount(id)
	if err != nil {
		http.Error(w, "unable to retrieve students count", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{Status: "success", Count: count})
	w.Header().Set("Content-Type", "application/json")
}
