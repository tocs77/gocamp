package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"rest-srv/models"
	"rest-srv/utility"
	"strconv"
	"strings"
)

func checkValidStudentField(field string) bool {
	validColumns := map[string]bool{
		"id":         true,
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
	}
	return validColumns[field]
}

func getStudentStructValues(model any) []any {
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

func GetStudentById(id int) (models.Student, error) {
	row := Db.QueryRow("SELECT * FROM students WHERE id = ?", id)
	var student models.Student
	err := row.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utility.ErrorHandler(err, "student not found")
	}
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "unable to retrieve student")
	}
	return student, nil
}

// GetStudents retrieves students with optional filters and sorting
// filters: map of field name to filter value (e.g., map[string]string{"email": "test@example.com"})
// sortParams: slice of strings in the format "field:asc" or "field:desc"
func GetStudents(filters map[string]string, sortParams []string) ([]models.Student, error) {
	var query string
	var orderByClauses []string
	var whereClauses []string
	var filterValues []any

	// Build WHERE clauses from filters
	for field, value := range filters {
		if !checkValidStudentField(field) {
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
		if !checkValidStudentField(field) {
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
	query = "SELECT * FROM students"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	if len(orderByClauses) > 0 {
		query += " ORDER BY " + strings.Join(orderByClauses, ", ")
	}

	studentsList := make([]models.Student, 0)
	var rows *sql.Rows
	var err error

	// Execute query with parameterized values
	if len(filterValues) > 0 {
		rows, err = Db.Query(query, filterValues...)
	} else {
		rows, err = Db.Query(query)
	}
	if err != nil {
		return nil, utility.ErrorHandler(err, "unable to retrieve students")
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utility.ErrorHandler(err, "unable to process student data")
		}
		studentsList = append(studentsList, student)
	}

	return studentsList, nil
}

func AddStudents(students []models.Student) ([]models.Student, error) {
	query := utility.GenerateInsertQuery(students[0], "students")
	stmt, err := Db.Prepare(query)
	if err != nil {
		return nil, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(students))
	for i, student := range students {
		res, err := stmt.Exec(getStudentStructValues(student)...)
		if err != nil {
			if strings.Contains(err.Error(), "Error 1452 (23000): Cannot add or update a child row: a foreign key constraint fails (`classes`.`students`, CONSTRAINT `1` FOREIGN KEY (`class`) REFERENCES `teachers` (`class`))") {
				return nil, utility.ErrorHandler(err, "class not found")
			}
			return nil, utility.ErrorHandler(err, "database error")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utility.ErrorHandler(err, "database error")
		}
		addedStudents[i] = student
		addedStudents[i].ID = int(lastID)
	}

	return addedStudents, nil
}

func PatchStudentFields(student *models.Student, updatedFields map[string]any) {
	studentVal := reflect.ValueOf(student).Elem()
	studentType := studentVal.Type()
	for key, value := range updatedFields {
		// Skip id field - it shouldn't be updated via PATCH
		if key == "id" {
			continue
		}
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			jsonTag := field.Tag.Get("json")
			// Extract the field name from the JSON tag (remove ",omitempty" if present)
			jsonFieldName := strings.Split(jsonTag, ",")[0]
			if jsonFieldName == key && studentVal.Field(i).CanSet() {
				studentVal.Field(i).Set(reflect.ValueOf(value).Convert(studentVal.Field(i).Type()))
				break
			}
		}
	}
}

func PatchStudent(id int, updateFields map[string]any) (models.Student, error) {
	student, err := GetStudentById(id)
	if err != nil {
		return models.Student{}, err
	}
	if student == (models.Student{}) {
		return models.Student{}, utility.ErrorHandler(errors.New("student not found"), "student not found")
	}

	// Apply patch updates to student
	PatchStudentFields(&student, updateFields)
	err = student.Validate()
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "invalid fields")
	}

	// Update database
	stmt, err := Db.Prepare("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?")
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(student.FirstName, student.LastName, student.Email, student.Class, student.ID)
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "database error")
	}
	return student, nil
}

func UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
	// Verify student exists before updating
	_, err := GetStudentById(id)
	if err != nil {
		return models.Student{}, err
	}

	// Set the ID from the parameter
	updatedStudent.ID = id

	// Update database
	stmt, err := Db.Prepare("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?")
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedStudent.FirstName, updatedStudent.LastName, updatedStudent.Email, updatedStudent.Class, id)
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "database error")
	}

	return updatedStudent, nil
}

func PatchStudents(updates []map[string]any) ([]models.Student, error) {
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

	updatedStudents := make([]models.Student, 0, len(updates))
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

		// Get existing student
		existingStudent, err := GetStudentById(id)
		if err != nil {
			rollbackNeeded = true
			return nil, err
		}

		// Apply patch updates to student
		PatchStudentFields(&existingStudent, update)

		// Update database within transaction
		stmt, err := tx.Prepare("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?")
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		_, err = stmt.Exec(existingStudent.FirstName, existingStudent.LastName, existingStudent.Email, existingStudent.Class, id)
		stmt.Close()
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		updatedStudents = append(updatedStudents, existingStudent)
	}

	// Commit the transaction if all updates succeeded
	err = tx.Commit()
	if err != nil {
		rollbackNeeded = true
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded = false

	return updatedStudents, nil
}

func DeleteStudent(id int) (models.Student, error) {
	// Get student before deleting to return it
	student, err := GetStudentById(id)
	if err != nil {
		return models.Student{}, err
	}

	stmt, err := Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Student{}, utility.ErrorHandler(err, "database error")
	}
	if rowsAffected == 0 {
		return models.Student{}, utility.ErrorHandler(errors.New("student not found"), "student not found")
	}

	return student, nil
}

func DeleteStudents(ids []int) ([]models.Student, error) {
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

	deletedStudents := make([]models.Student, 0, len(ids))
	for _, id := range ids {
		// Get student before deleting to return it
		student, err := GetStudentById(id)
		if err != nil {
			rollbackNeeded = true
			return nil, err
		}

		stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
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
			return nil, utility.ErrorHandler(errors.New("student not found"), "student not found")
		}

		deletedStudents = append(deletedStudents, student)
	}

	// Commit the transaction if all deletions succeeded
	err = tx.Commit()
	if err != nil {
		rollbackNeeded = true
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded = false

	return deletedStudents, nil
}
