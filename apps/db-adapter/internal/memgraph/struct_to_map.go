package memgraph

import (
	"reflect"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// StructToMap recursively converts a struct to a map using JSON tags.
// Nil pointers and unexported fields are excluded.
func StructToMap(input any) map[string]any {
	result := make(map[string]any)
	value := reflect.ValueOf(input)

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return result
		}
		value = value.Elem()
	}

	typ := value.Type()

	for i := range value.NumField() {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get the JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		// Only use the name before the first comma
		jsonKey := field.Name
		if jsonTag != "" {
			jsonKey = jsonTag
			if commaIdx := indexComma(jsonKey); commaIdx >= 0 {
				jsonKey = jsonKey[:commaIdx]
			}
		}

		// Skip empty json keys (e.g., tag is `json:"-"`)
		if jsonKey == "" {
			continue
		}

		// Handle nil pointers
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		// Dereference pointers
		val := fieldValue
		if fieldValue.Kind() == reflect.Ptr {
			val = fieldValue.Elem()
		}

		if isPreservedType(val.Interface()) {
			result[jsonKey] = val.Interface()

			continue
		}

		// Recurse into nested structs
		switch val.Kind() {
		case reflect.Struct:
			result[jsonKey] = StructToMap(val.Interface())
		case reflect.Slice, reflect.Array:
			result[jsonKey] = processSlice(val)
		default:
			result[jsonKey] = val.Interface()
		}
	}

	return result
}

func indexComma(tag string) int {
	for i, r := range tag {
		if r == ',' {
			return i
		}
	}
	return -1
}

// Checks if a value is one of the preserved types that shouldn't be expanded recursively
func isPreservedType(v any) bool {
	switch v.(type) {
	case dbtype.Point2D, *dbtype.Point2D,
		dbtype.Point3D, *dbtype.Point3D,
		time.Time,
		dbtype.LocalDateTime,
		dbtype.Date,
		dbtype.Time,
		dbtype.LocalTime,
		dbtype.Duration:
		return true
	default:
		return false
	}
}

func processSlice(val reflect.Value) []any {
	slice := make([]any, val.Len())
	for i := range val.Len() {
		item := val.Index(i).Interface()
		if isPreservedType(item) {
			slice[i] = item
		} else if reflect.ValueOf(item).Kind() == reflect.Struct {
			slice[i] = StructToMap(item)
		} else {
			slice[i] = item
		}
	}

	return slice
}
