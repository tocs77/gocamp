package utils

import (
	"reflect"
	"strings"
)

func MapStructFields(src any, dst any) {
	srcVal := reflect.ValueOf(src)
	if !srcVal.IsValid() {
		return
	}
	if srcVal.Kind() == reflect.Pointer {
		if srcVal.IsNil() {
			return
		}
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return
	}

	dstVal := reflect.ValueOf(dst)
	if !dstVal.IsValid() || dstVal.Kind() != reflect.Pointer || dstVal.IsNil() {
		return
	}
	dstVal = dstVal.Elem()
	if dstVal.Kind() != reflect.Struct {
		return
	}

	srcType := srcVal.Type()
	dstType := dstVal.Type()

	dstFieldByLowerName := make(map[string]reflect.Value, dstVal.NumField())
	for i := 0; i < dstVal.NumField(); i++ {
		fieldMeta := dstType.Field(i)
		field := dstVal.Field(i)
		if !field.CanSet() {
			continue
		}
		dstFieldByLowerName[strings.ToLower(fieldMeta.Name)] = field
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldMeta := srcType.Field(i)
		if !srcField.IsValid() {
			continue
		}

		dstField, ok := dstFieldByLowerName[strings.ToLower(srcFieldMeta.Name)]
		if !ok {
			continue
		}

		if srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
			continue
		}

		if srcField.Type().ConvertibleTo(dstField.Type()) {
			dstField.Set(srcField.Convert(dstField.Type()))
			continue
		}

		// Handle protobuf optional fields and similar pointer/value mismatches.
		if srcField.Kind() == reflect.Pointer && !srcField.IsNil() {
			elem := srcField.Elem()
			if elem.Type().AssignableTo(dstField.Type()) {
				dstField.Set(elem)
				continue
			}
			if elem.Type().ConvertibleTo(dstField.Type()) {
				dstField.Set(elem.Convert(dstField.Type()))
				continue
			}
		}

		if dstField.Kind() == reflect.Pointer {
			if srcField.Type().AssignableTo(dstField.Type().Elem()) {
				ptr := reflect.New(dstField.Type().Elem())
				ptr.Elem().Set(srcField)
				dstField.Set(ptr)
				continue
			}
			if srcField.Type().ConvertibleTo(dstField.Type().Elem()) {
				ptr := reflect.New(dstField.Type().Elem())
				ptr.Elem().Set(srcField.Convert(dstField.Type().Elem()))
				dstField.Set(ptr)
			}
		}
	}
}
