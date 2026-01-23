package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"rest-srv/db"
	"rest-srv/models"
	"rest-srv/utility"
)

func GetExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	exec, err := db.GetExecById(id)
	if err != nil {
		if err.Error() == "exec not found" {
			http.Error(w, "exec not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to retrieve exec", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(exec)
	w.Header().Set("Content-Type", "application/json")
}

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
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

	execsList, err := db.GetExecs(filters, sortByParams)
	if err != nil {
		http.Error(w, "unable to retrieve execs", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{Status: "success", Count: len(execsList), Data: execsList}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}

func AddExecHandler(w http.ResponseWriter, r *http.Request) {
	var newExecs []models.Exec
	err := json.NewDecoder(r.Body).Decode(&newExecs)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	for i := range newExecs {
		newExecs[i].UserCreatedAt = utility.NullString{NullString: sql.NullString{String: currentTime, Valid: true}}
		fmt.Println("Validating exec: ", newExecs[i])
		err = newExecs[i].Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if newExecs[i].Password == "" {
			http.Error(w, utility.ErrorHandler(errors.New("password is required"), "error adding exec").Error(), http.StatusBadRequest)
			return
		}

		encodedHash, err := utility.HashPassword(newExecs[i].Password)
		if err != nil {
			http.Error(w, utility.ErrorHandler(err, "error adding exec").Error(), http.StatusInternalServerError)
			return
		}
		newExecs[i].Password = encodedHash
	}

	addedExecs, err := db.AddExecs(newExecs)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(addedExecs)
	w.Header().Set("Content-Type", "application/json")
}

// execs/{id}
func UpdateExecHandler(w http.ResponseWriter, r *http.Request) {
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

	var updatedExec models.Exec
	err = json.NewDecoder(r.Body).Decode(&updatedExec)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err = updatedExec.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedExec, err = db.UpdateExec(id, updatedExec)
	if err != nil {
		if err.Error() == "exec not found" {
			http.Error(w, "exec not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update exec", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedExec)
	w.Header().Set("Content-Type", "application/json")
}

func PatchExecHandler(w http.ResponseWriter, r *http.Request) {
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

	updatedExec, err := db.PatchExec(id, updatedFields)
	if err != nil {
		if err.Error() == "exec not found" {
			http.Error(w, "exec not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update exec", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedExec)
	w.Header().Set("Content-Type", "application/json")
}

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedExecs, err := db.PatchExecs(updates)
	if err != nil {
		if err.Error() == "exec not found" || strings.Contains(err.Error(), "exec not found") {
			http.Error(w, "exec not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to update execs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedExecs)
	w.Header().Set("Content-Type", "application/json")
}

func DeleteExecHandler(w http.ResponseWriter, r *http.Request) {
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

	deletedExec, err := db.DeleteExec(id)
	if err != nil {
		if err.Error() == "exec not found" {
			http.Error(w, "exec not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to delete exec", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		ID      int    `json:"id"`
	}{Status: "success", Message: "Exec deleted successfully", ID: deletedExec.ID}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteExecsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deletedExecs, err := db.DeleteExecs(ids)
	if err != nil {
		if strings.Contains(err.Error(), "exec not found") {
			http.Error(w, "exec not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unable to delete execs", http.StatusInternalServerError)
		return
	}

	deletedIDs := make([]int, len(deletedExecs))
	for i, exec := range deletedExecs {
		deletedIDs[i] = exec.ID
	}

	json.NewEncoder(w).Encode(struct {
		Status     string `json:"status"`
		Message    string `json:"message"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{Status: "success", Message: "Execs deleted successfully", DeletedIDs: deletedIDs})
	w.Header().Set("Content-Type", "application/json")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Data validation
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Search for exec by username
	exec, err := db.GetExecByUsername(loginData.Username)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	// is user active?
	if exec.InactiveStatus {
		http.Error(w, "user is inactive", http.StatusUnauthorized)
		return
	}

	// verify password
	valid, err := utility.ComparePassword(exec.Password, loginData.Password)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	if !valid {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	// generate token
	testToken, err := utility.SingToken(strconv.Itoa(exec.ID), exec.Username, exec.Role)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// return token
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{Name: "Bearer", Value: testToken, Path: "/", HttpOnly: true, Secure: true, Expires: time.Now().Add(24 * time.Hour)})
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	}{Status: "success", Token: testToken})
}
