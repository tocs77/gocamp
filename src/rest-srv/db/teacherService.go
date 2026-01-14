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

func checkValidField(field string) bool {
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

func getStructValues(model any) []any {
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

func GetTeacherById(id int) (models.Teacher, error) {
	row := Db.QueryRow("SELECT * FROM teachers WHERE id = ?", id)
	var teacher models.Teacher
	err := row.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utility.ErrorHandler(err, "teacher not found")
	}
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "unable to retrieve teacher")
	}
	return teacher, nil
}

// GetTeachers retrieves teachers with optional filters and sorting
// filters: map of field name to filter value (e.g., map[string]string{"email": "test@example.com"})
// sortParams: slice of strings in the format "field:asc" or "field:desc"
func GetTeachers(filters map[string]string, sortParams []string) ([]models.Teacher, error) {
	var query string
	var orderByClauses []string
	var whereClauses []string
	var filterValues []any

	// Build WHERE clauses from filters
	for field, value := range filters {
		if !checkValidField(field) {
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
		rows, err = Db.Query(query, filterValues...)
	} else {
		rows, err = Db.Query(query)
	}
	if err != nil {
		return nil, utility.ErrorHandler(err, "unable to retrieve teachers")
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, utility.ErrorHandler(err, "unable to process teacher data")
		}
		teachersList = append(teachersList, teacher)
	}

	return teachersList, nil
}

func AddTeachers(teachers []models.Teacher) ([]models.Teacher, error) {
	query := utility.GenerateInsertQuery(teachers[0], "teachers")
	stmt, err := Db.Prepare(query)
	if err != nil {
		return nil, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(teachers))
	for i, teacher := range teachers {
		res, err := stmt.Exec(getStructValues(teacher)...)
		if err != nil {
			return nil, utility.ErrorHandler(err, "database error")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utility.ErrorHandler(err, "database error")
		}
		addedTeachers[i] = teacher
		addedTeachers[i].ID = int(lastID)
	}

	return addedTeachers, nil
}

func PatchTeacherFields(teacher *models.Teacher, updatedFields map[string]any) {
	teacherVal := reflect.ValueOf(teacher).Elem()
	teacherType := teacherVal.Type()
	for key, value := range updatedFields {
		// Skip id field - it shouldn't be updated via PATCH
		if key == "id" {
			continue
		}
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
}

func PatchTeacher(id int, updateFields map[string]any) (models.Teacher, error) {
	teacher, err := GetTeacherById(id)
	if err != nil {
		return models.Teacher{}, err
	}
	if teacher == (models.Teacher{}) {
		return models.Teacher{}, utility.ErrorHandler(errors.New("teacher not found"), "teacher not found")
	}

	// Apply patch updates to teacher
	PatchTeacherFields(&teacher, updateFields)
	err = teacher.Validate()
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "invalid fields")
	}

	// Update database
	stmt, err := Db.Prepare("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?")
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID)
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "database error")
	}
	return teacher, nil
}

func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	// Verify teacher exists before updating
	_, err := GetTeacherById(id)
	if err != nil {
		return models.Teacher{}, err
	}

	// Set the ID from the parameter
	updatedTeacher.ID = id

	// Update database
	stmt, err := Db.Prepare("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?")
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, id)
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "database error")
	}

	return updatedTeacher, nil
}
func PatchTeachers(updates []map[string]any) ([]models.Teacher, error) {
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

	updatedTeachers := make([]models.Teacher, 0, len(updates))
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

		// Get existing teacher
		existingTeacher, err := GetTeacherById(id)
		if err != nil {
			rollbackNeeded = true
			return nil, err
		}

		// Apply patch updates to teacher
		PatchTeacherFields(&existingTeacher, update)

		// Update database within transaction
		stmt, err := tx.Prepare("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?")
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		_, err = stmt.Exec(existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, id)
		stmt.Close()
		if err != nil {
			rollbackNeeded = true
			return nil, utility.ErrorHandler(err, "database error")
		}
		updatedTeachers = append(updatedTeachers, existingTeacher)
	}

	// Commit the transaction if all updates succeeded
	err = tx.Commit()
	if err != nil {
		rollbackNeeded = true
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded = false

	return updatedTeachers, nil
}

func DeleteTeacher(id int) (models.Teacher, error) {
	// Get teacher before deleting to return it
	teacher, err := GetTeacherById(id)
	if err != nil {
		return models.Teacher{}, err
	}

	stmt, err := Db.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "database error")
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "database error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Teacher{}, utility.ErrorHandler(err, "database error")
	}
	if rowsAffected == 0 {
		return models.Teacher{}, utility.ErrorHandler(errors.New("teacher not found"), "teacher not found")
	}

	return teacher, nil
}

func DeleteTeachers(ids []int) ([]models.Teacher, error) {
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

	deletedTeachers := make([]models.Teacher, 0, len(ids))
	for _, id := range ids {
		// Get teacher before deleting to return it
		teacher, err := GetTeacherById(id)
		if err != nil {
			rollbackNeeded = true
			return nil, err
		}

		stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
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
			return nil, utility.ErrorHandler(errors.New("teacher not found"), "teacher not found")
		}

		deletedTeachers = append(deletedTeachers, teacher)
	}

	// Commit the transaction if all deletions succeeded
	err = tx.Commit()
	if err != nil {
		rollbackNeeded = true
		return nil, utility.ErrorHandler(err, "database error")
	}
	rollbackNeeded = false

	return deletedTeachers, nil
}

func GetTeacherStudents(id int) ([]models.Student, error) {
	rows, err := Db.Query("SELECT * FROM students WHERE class = (SELECT class FROM teachers WHERE id = ?)", id)
	if err != nil {
		return nil, utility.ErrorHandler(err, "database error")
	}
	defer rows.Close()

	students := make([]models.Student, 0)
	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utility.ErrorHandler(err, "unable to process student data")
		}
		students = append(students, student)
	}

	return students, nil
}

func GetTeacherStudentsCount(id int) (int, error) {
	rows, err := Db.Query("SELECT COUNT(*) FROM students WHERE class = (SELECT class FROM teachers WHERE id = ?)", id)
	if err != nil {
		return 0, utility.ErrorHandler(err, "database error")
	}
	defer rows.Close()

	if rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			return 0, utility.ErrorHandler(err, "unable to process student data")
		}
		return count, nil
	}
	return 0, nil
}
