package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if filtersVal.Kind() == reflect.Pointer {
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
		if filterVal.Kind() == reflect.Pointer {
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
				fmt.Println("columnName", columnName)
				if columnName == "_id" {
					id, err := primitive.ObjectIDFromHex(field.String())
					if err != nil {
						return nil, fmt.Errorf("invalid ID: %w", err)
					}
					resultFilter[columnName] = id
					continue
				}
				if field.Kind() == reflect.String {
					quoted := regexp.QuoteMeta(field.String())
					resultFilter[columnName] = bson.M{
						"$regex":   "^" + quoted + "$",
						"$options": "i",
					}
					continue
				}
				resultFilter[columnName] = field.Interface()
			}
		}
	}
	return resultFilter, nil
}

func BuildSortForModel(model any, sorts any, descValue any) (bson.D, error) {
	modelType := reflect.TypeOf(model)
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct or pointer to struct")
	}

	allowedSortFields := make(map[string]struct{}, modelType.NumField())
	for i := 0; i < modelType.NumField(); i++ {
		bsonTag := modelType.Field(i).Tag.Get("bson")
		columnName := strings.Split(bsonTag, ",")[0]
		if columnName != "" {
			allowedSortFields[columnName] = struct{}{}
		}
	}

	sortOptions := bson.D{}
	sortsVal := reflect.ValueOf(sorts)
	if !sortsVal.IsValid() {
		return sortOptions, nil
	}
	if sortsVal.Kind() == reflect.Pointer {
		if sortsVal.IsNil() {
			return sortOptions, nil
		}
		sortsVal = sortsVal.Elem()
	}

	var sortItems []reflect.Value
	switch sortsVal.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < sortsVal.Len(); i++ {
			sortItems = append(sortItems, sortsVal.Index(i))
		}
	case reflect.Struct:
		sortItems = append(sortItems, sortsVal)
	default:
		return nil, fmt.Errorf("sorts must be a struct, pointer to struct, slice or array")
	}

	for _, sortVal := range sortItems {
		if sortVal.Kind() == reflect.Interface {
			if sortVal.IsNil() {
				continue
			}
			sortVal = sortVal.Elem()
		}
		if sortVal.Kind() == reflect.Pointer {
			if sortVal.IsNil() {
				continue
			}
			sortVal = sortVal.Elem()
		}

		fieldName, orderVal, ok := extractSortParts(sortVal)
		if !ok {
			continue
		}
		if _, exists := allowedSortFields[fieldName]; !exists {
			continue
		}

		order := 1
		if isSortDesc(orderVal, descValue) {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{Key: fieldName, Value: order})
	}

	return sortOptions, nil
}

func columnNameFromFilterField(field reflect.StructField) string {
	if bsonTag := strings.Split(field.Tag.Get("bson"), ",")[0]; bsonTag != "" {
		return bsonTag
	}

	if protobufName := protobufFieldName(field.Tag.Get("protobuf")); protobufName != "" {
		if protobufName == "id" {
			return "_id"
		}
		return protobufName
	}

	jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
	if jsonName == "id" {
		return "_id"
	}
	return jsonName
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

func extractSortParts(sortVal reflect.Value) (string, reflect.Value, bool) {
	fieldMethod := sortVal.MethodByName("GetField")
	orderMethod := sortVal.MethodByName("GetOrder")
	if fieldMethod.IsValid() && orderMethod.IsValid() {
		fieldResults := fieldMethod.Call(nil)
		orderResults := orderMethod.Call(nil)
		if len(fieldResults) == 1 && len(orderResults) == 1 {
			fieldResult := fieldResults[0]
			if fieldResult.Kind() == reflect.String {
				return fieldResult.String(), orderResults[0], true
			}
		}
	}

	if sortVal.Kind() != reflect.Struct {
		return "", reflect.Value{}, false
	}

	field := sortVal.FieldByName("Field")
	order := sortVal.FieldByName("Order")
	if !field.IsValid() || !order.IsValid() {
		return "", reflect.Value{}, false
	}
	if field.Kind() != reflect.String || !field.CanInterface() {
		return "", reflect.Value{}, false
	}

	return field.Interface().(string), order, true
}

func isSortDesc(orderVal reflect.Value, descValue any) bool {
	if !orderVal.IsValid() {
		return false
	}

	for orderVal.Kind() == reflect.Interface || orderVal.Kind() == reflect.Pointer {
		if orderVal.IsNil() {
			return false
		}
		orderVal = orderVal.Elem()
	}

	if descValue != nil && orderVal.CanInterface() && reflect.DeepEqual(orderVal.Interface(), descValue) {
		return true
	}

	if orderVal.Kind() == reflect.String {
		return strings.EqualFold(orderVal.String(), "DESC")
	}

	switch orderVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return orderVal.Int() == -1
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return orderVal.Uint() == 1 && descValue == nil
	}

	return false
}
