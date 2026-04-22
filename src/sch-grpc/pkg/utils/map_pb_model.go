package utils

import (
	"reflect"
)

func MapPBModelToModel(src any, dst any) {
	srcVal := reflect.ValueOf(src)
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

	dstVal := reflect.ValueOf(dst)
	if !dstVal.IsValid() || dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return
	}
	dstVal = dstVal.Elem()
	if dstVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < dstVal.NumField(); i++ {
		dstField := dstVal.Field(i)
		if !dstField.CanSet() {
			continue
		}

		dstFieldMeta := dstVal.Type().Field(i)
		srcField := srcVal.FieldByName(dstFieldMeta.Name)
		if !srcField.IsValid() {
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