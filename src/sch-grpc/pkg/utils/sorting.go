package utils

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildFilterForModel(model any, filters any) (bson.M, error) {
	modelType := reflect.TypeOf(model)
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct or pointer to struct")
	}
	allowedFilterFields := make(map[string]struct{}, modelType.NumField())
	resultFilter := bson.M{}

	for i := 0; i < modelType.NumField(); i++ {
		bsonTag := modelType.Field(i).Tag.Get("bson")
		columnName := strings.Split(bsonTag, ",")[0]
		if columnName != "" {
			allowedFilterFields[columnName] = struct{}{}
		}
	}

	filtersVal := reflect.ValueOf(filters)
	if !filtersVal.IsValid() {
		return resultFilter, nil
	}
	if filtersVal.Kind() == reflect.Ptr {
		if filtersVal.IsNil() {
			return resultFilter, nil
		}
		filtersVal = filtersVal.Elem()
	}

	var filterItems []reflect.Value
	switch filtersVal.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < filtersVal.Len(); i++ {
			filterItems = append(filterItems, filtersVal.Index(i))
		}
	case reflect.Struct:
		filterItems = append(filterItems, filtersVal)
	default:
		return nil, fmt.Errorf("filters must be a struct, pointer to struct, slice or array")
	}

	for _, filterVal := range filterItems {
		if filterVal.Kind() == reflect.Interface {
			if filterVal.IsNil() {
				continue
			}
			filterVal = filterVal.Elem()
		}
		if filterVal.Kind() == reflect.Ptr {
			if filterVal.IsNil() {
				continue
			}
			filterVal = filterVal.Elem()
		}
		if filterVal.Kind() != reflect.Struct {
			return nil, fmt.Errorf("each filter must be a struct or pointer to struct")
		}

		filterType := filterVal.Type()
		for i := 0; i < filterVal.NumField(); i++ {
			field := filterVal.Field(i)
			fieldMeta := filterType.Field(i)
			if fieldMeta.PkgPath != "" {
				// Unexported fields (e.g. protobuf internal state) are not readable via Interface().
				continue
			}
			if field.IsValid() && !field.IsZero() {
				columnName := columnNameFromFilterField(fieldMeta)
				if columnName == "" {
					continue
				}
				if _, ok := allowedFilterFields[columnName]; !ok {
					continue
				}
				if !field.CanInterface() {
					continue
				}
				resultFilter[columnName] = field.Interface()
			}
		}
	}
	return resultFilter, nil
}

func columnNameFromFilterField(field reflect.StructField) string {
	if bsonTag := strings.Split(field.Tag.Get("bson"), ",")[0]; bsonTag != "" {
		return bsonTag
	}

	if protobufName := protobufFieldName(field.Tag.Get("protobuf")); protobufName != "" {
		return protobufName
	}

	return strings.Split(field.Tag.Get("json"), ",")[0]
}

func protobufFieldName(tag string) string {
	if tag == "" {
		return ""
	}

	for part := range strings.SplitSeq(tag, ",") {
		if name, ok := strings.CutPrefix(part, "name="); ok {
			return name
		}
	}

	return ""
}
