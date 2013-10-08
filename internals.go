package mantra

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	uppercaseLetters = regexp.MustCompile("([A-Z])")
)

func orderedStructFields(object interface{}) []string {
	s := reflect.ValueOf(object).Elem()
	typeOfT := s.Type()

	fieldNames := make([]string, 0, s.NumField())

	for i := 0; i < s.NumField(); i++ {
		field := typeOfT.Field(i)

		if fieldTag := field.Tag.Get("mantraIgnore"); len(fieldTag) == 0 {
			fieldNames = append(fieldNames, typeOfT.Field(i).Name)
		}
	}

	return fieldNames
}

func orderedDatabaseFields(object interface{}) []string {
	fieldNames := orderedStructFields(object)
	fieldMap := structToDatabaseFieldMap(object)

	orderedColumnNames := make([]string, 0, len(fieldNames))
	for _, name := range fieldNames {
		orderedColumnNames = append(orderedColumnNames, fieldMap[name])
	}

	return orderedColumnNames
}

func structToDatabaseFieldMap(object interface{}) map[string]string {
	s := reflect.ValueOf(object).Elem()
	typeOfT := s.Type()

	fieldNames := make(map[string]string, s.NumField())

	for i := 0; i < s.NumField(); i++ {
		field := typeOfT.Field(i)
		// If there's a 'column' tag on this field, use that as the column name
		if fieldTag := field.Tag.Get("mantraColumn"); len(fieldTag) > 0 {
			fieldNames[field.Name] = fieldTag
		} else {
			// if there was no 'column' tag, convert the field name to underscore_case and use that
			fieldNames[field.Name] = strings.ToLower(uppercaseLetters.ReplaceAllString(field.Name, "_$1")[1:])
		}
	}

	return fieldNames
}

func fieldWithTag(object interface{}, tagName string) string {
	s := reflect.ValueOf(object).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		field := typeOfT.Field(i)
		if fieldTag := field.Tag.Get(tagName); len(fieldTag) > 0 {
			mapping := structToDatabaseFieldMap(object)
			return mapping[field.Name]
		}
	}

	return ""
}

func insertFieldList(object interface{}) string {
	orderedColumnNames := orderedDatabaseFields(object)
	return fmt.Sprintf("`%s`", strings.Join(orderedColumnNames, "`, `"))
}

func updateFieldList(object interface{}) string {
	fieldNames := orderedDatabaseFields(object)
	values := make([]string, 0, len(fieldNames))

	for _, field := range fieldNames {
		values = append(values, fmt.Sprintf("%v = ?", field))
	}

	return strings.Join(values, ", ")
}

func primaryKeyField(object interface{}) string {
	// First, check for a field tagged with mantraPrimaryKey
	if field := fieldWithTag(object, "mantraPrimaryKey"); len(field) > 0 {
		return field
	}

	// If there's no tagged field, then fall back to looking for a "id" field
	fieldMap := structToDatabaseFieldMap(object)
	if _, ok := fieldMap["Id"]; ok {
		return "id"
	}

	// If we didn't find anything, return an empty string
	return ""
}

func deletionFlag(object interface{}) string {
	// First, check for a field tagged with mantraDeletionFlag
	if field := fieldWithTag(object, "mantraDeletionFlag"); len(field) > 0 {
		return field
	}

	// If there's no tagged field, then fall back to looking for a "deleted" field
	fieldMap := structToDatabaseFieldMap(object)
	if _, ok := fieldMap["Deleted"]; ok {
		return "deleted"
	}

	// If we didn't find anything, return an empty string
	return ""
}

func valuePlaceholders(object interface{}) string {
	s := reflect.ValueOf(object).Elem()

	fieldNames := make([]string, 0, s.NumField())

	for i := 0; i < s.NumField(); i++ {
		fieldNames = append(fieldNames, "?")
	}

	return strings.Join(fieldNames, " ")
}
