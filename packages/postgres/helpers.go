package postgres

import (
	"reflect"
	"strconv"
	"strings"
)

func taggedFieldsFromStruct(object interface{}, tag string) []string {
	fields := []string{}

	obj := reflect.TypeOf(object)
	if obj.Kind() != reflect.Struct {
		panic("given type is not a struct")
	}

	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		if name, ok := field.Tag.Lookup(tag); ok && field.PkgPath == "" {
			fields = append(fields, name)
		}
	}

	return fields
}

func valuesFromTaggedFields(object interface{}, tag string) []interface{} {
	values := []interface{}{}

	obj := reflect.ValueOf(object)
	objType := obj.Type()
	if objType.Kind() != reflect.Struct {
		panic("given type is not a struct")
	}

	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		fieldType := objType.Field(i)
		if _, ok := fieldType.Tag.Lookup(tag); ok && fieldType.PkgPath == "" {
			values = append(values, field.Interface())
		}
	}

	return values
}

//-------------------------------------------------------------------------------------------------

func valuePlaceholders(count, length int) string {
	placeholders := make([]string, 0, count)
	for i := 0; i < count; i++ {
		placeholders = append(placeholders, valuePlaceholder(i*length+1, length))
	}
	return strings.Join(placeholders, ",")
}

func valuePlaceholder(offset, length int) string {
	var sb strings.Builder
	sb.WriteByte(40) // "("
	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteByte(44) // ","
		}
		sb.WriteByte(36) // "$"
		sb.WriteString(strconv.Itoa(offset + i))
	}
	sb.WriteByte(41) // ")"
	return sb.String()
}
