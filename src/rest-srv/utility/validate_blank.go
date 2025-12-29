package utility

import (
	"fmt"
	"reflect"
	"strings"
)

func ValidateBlank(value any) error {
	val := reflect.ValueOf(value).Elem()
	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := valType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}
		columnName := strings.Split(dbTag, ",")[0]
		if columnName == "id" {
			continue
		}
		if strings.Contains(dbTag, "not_null") {
			if val.Field(i).IsZero() {
				return fmt.Errorf("field %s is required", columnName)
			}
		}
	}
	return nil
}
