package main

import (
	"reflect"
	"strings"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func reflectValue(obj interface{}) reflect.Value {
	var val reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		val = reflect.ValueOf(obj).Elem()
	} else {
		val = reflect.ValueOf(obj)
	}

	return val
}

func isExportableField(field reflect.StructField) bool {
	// PkgPath is empty for exported fields.
	return field.PkgPath == ""
}

// as struct tags can continue multiple values, this gets the first
func GetFirstTagValue(tag reflect.StructTag, attr string) (value string) {
	tagValue := tag.Get(attr)
	if len(tagValue) == 0 {
		return
	}
	value = strings.Split(tagValue, ",")[0]
	return
}

// get field name based on struct tags
func GetFieldName(field reflect.StructField) (value string) {
	value = field.Tag.Get("access_field_name")
	if len(value) > 0 {
		return
	}
	value = GetFirstTagValue(field.Tag, "json")
	if len(value) > 0 {
		return
	}
	value = field.Name
	return
}

func GetChangedFields(a1, a2 interface{}) (err error, changes []string) {
	v1 := reflectValue(a1)
	v2 := reflectValue(a2)
	objType := v2.Type()
	fieldsCount := v2.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		val1 := v1.Field(i)
		val2 := v2.Field(i)
		i1, i2 := val1.Interface(), val2.Interface()
		if isExportableField(field) {
			switch i1.(type) {
			case int, bool, string, float64:
				if i1 != i2 {
					fieldName := GetFieldName(field)
					changes = append(changes, fieldName)
				}
			default:
				// TODO handle unsupported types
			}
		}
	}
	return
}
