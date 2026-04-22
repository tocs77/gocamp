package utils

import (
	"reflect"
	"strings"
)

func MapModelToPB(model any, pb any) {
	srcVal := reflect.ValueOf(model)
	if !srcVal.IsValid() {
		return
	}
	if srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return
		}
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return
	}

	dstVal := reflect.ValueOf(pb)
	if !dstVal.IsValid() || dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return
	}
	dstVal = dstVal.Elem()
	if dstVal.Kind() != reflect.Struct {
		return
	}

	srcType := srcVal.Type()
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldMeta := srcType.Field(i)

		if !srcField.IsValid() {
			continue
		}

		dstField := findFieldByName(dstVal, srcFieldMeta.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		if srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
			continue
		}

		if srcField.Type().ConvertibleTo(dstField.Type()) {
			dstField.Set(srcField.Convert(dstField.Type()))
		}
	}
}

func findFieldByName(dstVal reflect.Value, sourceName string) reflect.Value {
	if field := dstVal.FieldByName(sourceName); field.IsValid() {
		return field
	}

	// protobuf-generated names commonly use "Id", while models often use "ID".
	if strings.Contains(sourceName, "ID") {
		protoName := strings.ReplaceAll(sourceName, "ID", "Id")
		if field := dstVal.FieldByName(protoName); field.IsValid() {
			return field
		}
	}

	return reflect.Value{}
}
