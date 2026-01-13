package utility

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateInsertQuery(model any, tableName string) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
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

		if columns != "" {
			columns += ", "
			placeholders += ", "
		}
		columns += columnName
		placeholders += "?"
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}
