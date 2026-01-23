package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"rest-srv/models"
	"rest-srv/utility"
)

func checkValidExecField(field string) bool {
	validColumns := map[string]bool{
		"id":                   true,
		"first_name":           true,
		"last_name":            true,
		"email":                true,
		"username":             true,
		"password":             true,
		"password_changed_at":  true,
		"user_created_at":      true,
		"password_reset_token": true,
		"inactive_status":      true,
		"role":                 true,
	}
	return validColumns[field]
}

func getExecStructValues(model any) []any {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)
	var values []any
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}
		// Extract the column name from the db tag (remove metadata like "not_null", "unique", etc.)
		columnName := strings.Split(dbTag, ",")[0]
		if columnName == "id" {
			continue
		}
		values = append(values, modelValue.Field(i).Interface())
	}
	return values
}

func GetExecById(id int) (models.Exec, error) {
	row := Db.QueryRow("SELECT id, first_name, last_name, email, username, password, password_changed_at, user_created_at, password_reset_token, inactive_status, role FROM execs WHERE id = ?", id)
	var exec models.Exec
	err := row.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.PasswordChangedAt, &exec.UserCreatedAt, &exec.PasswordResetToken, &exec.InactiveStatus, &exec.Role)
	if err == sql.ErrNoRows {
		return models.Exec{}, utility.ErrorHandler(err, "exec not found")
	}
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "unable to retrieve exec")
	}
	return exec, nil
}

func GetExecByUsername(username string) (models.Exec, error) {
	row := Db.QueryRow("SELECT id, first_name, last_name, email, username, password, password_changed_at, user_created_at, password_reset_token, inactive_status, role FROM execs WHERE username = ?", username)
	var exec models.Exec
	err := row.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.PasswordChangedAt, &exec.UserCreatedAt, &exec.PasswordResetToken, &exec.InactiveStatus, &exec.Role)
	if err == sql.ErrNoRows {
		return models.Exec{}, utility.ErrorHandler(err, "exec not found")
	}
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "unable to retrieve exec")
	}
	return exec, nil
}

// GetExecs retrieves execs with optional filters and sorting
// filters: map of field name to filter value (e.g., map[string]string{"email": "test@example.com"})
// sortParams: slice of strings in the format "field:asc" or "field:desc"
func GetExecs(filters map[string]string, sortParams []string) ([]models.Exec, error) {
	var query string
	var orderByClauses []string
	var whereClauses []string
	var filterValues []any

	// Build WHERE clauses from filters
	for field, value := range filters {
		if !checkValidExecField(field) {
			continue // Skip invalid field
		}
		if value != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			filterValues = append(filterValues, value)
		}
	}

	// Parse each sortBy parameter (format: field:asc or field:desc)
	for _, sortParam := range sortParams {
		parts := strings.Split(sortParam, ":")
		if len(parts) != 2 {
			continue // Skip invalid format
		}

		field := strings.TrimSpace(parts[0])
		order := strings.TrimSpace(strings.ToUpper(parts[1]))

		// Validate field
		if !checkValidExecField(field) {
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
	query = "SELECT id, first_name, last_name, email, username, password, password_changed_at, user_created_at, password_reset_token, inactive_status, role FROM execs"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	if len(orderByClauses) > 0 {
		query += " ORDER BY " + strings.Join(orderByClauses, ", ")
	}

	execsList := make([]models.Exec, 0)
	var rows *sql.Rows
	var err error

	// Execute query with parameterized values
	if len(filterValues) > 0 {
		rows, err = Db.Query(query, filterValues...)
	} else {
		rows, err = Db.Query(query)
	}
	if err != nil {
		return nil, utility.ErrorHandler(err, "unable to retrieve execs")
	}
	defer rows.Close()

	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.PasswordChangedAt, &exec.UserCreatedAt, &exec.PasswordResetToken, &exec.InactiveStatus, &exec.Role)
		if err != nil {
			return nil, utility.ErrorHandler(err, "unable to process exec data")
		}
		execsList = append(execsList, exec)
	}

	return execsList, nil
}

func AddExecs(execs []models.Exec) ([]models.Exec, error) {
	query := utility.GenerateInsertQuery(execs[0], "execs")
	stmt, err := Db.Prepare(query)
	if err != nil {
		return nil, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(execs))
	for i, exec := range execs {
		res, err := stmt.Exec(getExecStructValues(exec)...)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				if strings.Contains(err.Error(), "email") {
					return nil, utility.ErrorHandler(err, "email already exists")
				}
				if strings.Contains(err.Error(), "username") {
					return nil, utility.ErrorHandler(err, "username already exists")
				}
			}
			return nil, utility.ErrorHandler(err, "database error")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utility.ErrorHandler(err, "database error")
		}
		addedExecs[i] = exec
		addedExecs[i].ID = int(lastID)
	}

	return addedExecs, nil
}

func PatchExecFields(exec *models.Exec, updatedFields map[string]any) {
	execVal := reflect.ValueOf(exec).Elem()
	execType := execVal.Type()
	for key, value := range updatedFields {
		// Skip id field - it shouldn't be updated via PATCH
		if key == "id" {
			continue
		}
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)
			jsonTag := field.Tag.Get("json")
			// Extract the field name from the JSON tag (remove ",omitempty" if present)
			jsonFieldName := strings.Split(jsonTag, ",")[0]
			if jsonFieldName == key && execVal.Field(i).CanSet() {
				execVal.Field(i).Set(reflect.ValueOf(value).Convert(execVal.Field(i).Type()))
				break
			}
		}
	}
}

func PatchExec(id int, updateFields map[string]any) (models.Exec, error) {
	exec, err := GetExecById(id)
	if err != nil {
		return models.Exec{}, err
	}
	if exec == (models.Exec{}) {
		return models.Exec{}, utility.ErrorHandler(errors.New("exec not found"), "exec not found")
	}

	// Apply patch updates to exec
	PatchExecFields(&exec, updateFields)
	err = exec.Validate()
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "invalid fields")
	}

	// Update database
	stmt, err := Db.Prepare("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, password_changed_at = ?, user_created_at = ?, password_reset_token = ?, inactive_status = ?, role = ? WHERE id = ?")
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(exec.FirstName, exec.LastName, exec.Email, exec.Username, exec.Password, exec.PasswordChangedAt, exec.UserCreatedAt, exec.PasswordResetToken, exec.InactiveStatus, exec.Role, exec.ID)
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "database error")
	}
	return exec, nil
}

func UpdateExec(id int, updatedExec models.Exec) (models.Exec, error) {
	// Verify exec exists before updating
	_, err := GetExecById(id)
	if err != nil {
		return models.Exec{}, err
	}

	// Set the ID from the parameter
	updatedExec.ID = id

	// Update database
	stmt, err := Db.Prepare("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, password_changed_at = ?, user_created_at = ?, password_reset_token = ?, inactive_status = ?, role = ? WHERE id = ?")
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedExec.FirstName, updatedExec.LastName, updatedExec.Email, updatedExec.Username, updatedExec.Password, updatedExec.PasswordChangedAt, updatedExec.UserCreatedAt, updatedExec.PasswordResetToken, updatedExec.InactiveStatus, updatedExec.Role, id)
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "database error")
	}

	return updatedExec, nil
}

func PatchExecs(updates []map[string]any) ([]models.Exec, error) {
	tx, err := Db.Begin()
	if err != nil {
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded := true
	defer func() {
		if rollbackNeeded {
			tx.Rollback()
		}
	}()

	updatedExecs := make([]models.Exec, 0, len(updates))
	for _, update := range updates {
		// Extract and validate ID
		idVal, ok := update["id"]
		if !ok {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(errors.New("id is required"), "id is required")
		}

		var id int
		switch v := idVal.(type) {
		case string:
			var err error
			id, err = strconv.Atoi(v)
			if err != nil {
				rollbackNeeded = true
				return nil, utility.ErrorHandler(err, "invalid id")
			}
		case float64:
			id = int(v)
		case int:
			id = v
		default:
			rollbackNeeded = true
			return nil, utility.ErrorHandler(errors.New("invalid id type"), "invalid id type")
		}

		if id == 0 {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(errors.New("no id"), "no id")
		}

		// Get existing exec
		existingExec, err := GetExecById(id)
		if err != nil {
			rollbackNeeded = true
			return nil, err
		}

		// Apply patch updates to exec
		PatchExecFields(&existingExec, update)

		// Update database within transaction
		stmt, err := tx.Prepare("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, password_changed_at = ?, user_created_at = ?, password_reset_token = ?, inactive_status = ?, role = ? WHERE id = ?")
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		_, err = stmt.Exec(existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.Password, existingExec.PasswordChangedAt, existingExec.UserCreatedAt, existingExec.PasswordResetToken, existingExec.InactiveStatus, existingExec.Role, id)
		stmt.Close()
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		updatedExecs = append(updatedExecs, existingExec)
	}

	// Commit the transaction if all updates succeeded
	err = tx.Commit()
	if err != nil {
		rollbackNeeded = true
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded = false

	return updatedExecs, nil
}

func DeleteExec(id int) (models.Exec, error) {
	// Get exec before deleting to return it
	exec, err := GetExecById(id)
	if err != nil {
		return models.Exec{}, err
	}

	stmt, err := Db.Prepare("DELETE FROM execs WHERE id = ?")
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Exec{}, utility.ErrorHandler(err, "database error")
	}
	if rowsAffected == 0 {
		return models.Exec{}, utility.ErrorHandler(errors.New("exec not found"), "exec not found")
	}

	return exec, nil
}

func DeleteExecs(ids []int) ([]models.Exec, error) {
	tx, err := Db.Begin()
	if err != nil {
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded := true
	defer func() {
		if rollbackNeeded {
			tx.Rollback()
		}
	}()

	deletedExecs := make([]models.Exec, 0, len(ids))
	for _, id := range ids {
		// Get exec before deleting to return it
		exec, err := GetExecById(id)
		if err != nil {
			rollbackNeeded = true
			return nil, err
		}

		stmt, err := tx.Prepare("DELETE FROM execs WHERE id = ?")
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		result, err := stmt.Exec(id)
		stmt.Close()
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		if rowsAffected == 0 {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(errors.New("exec not found"), "exec not found")
		}

		deletedExecs = append(deletedExecs, exec)
	}

	// Commit the transaction if all deletions succeeded
	err = tx.Commit()
	if err != nil {
		rollbackNeeded = true
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded = false

	return deletedExecs, nil
}
