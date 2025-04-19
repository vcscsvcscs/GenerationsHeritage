package memgraph

import (
	"reflect"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type NestedStruct struct {
	NestedField string `json:"nested_field"`
}

type TestStruct struct {
	ExportedField   string       `json:"exported_field"`
	unexportedField string       // Should be ignored
	IgnoredField    string       `json:"-"`
	PointerField    *string      `json:"pointer_field"`
	Nested          NestedStruct `json:"nested"`
	NilPointer      *string      `json:"nil_pointer"`
}

func TestStructToMap(t *testing.T) {
	// Test data
	pointerValue := "pointer value"
	testStruct := TestStruct{
		ExportedField:   "exported value",
		unexportedField: "unexported value",
		IgnoredField:    "ignored value",
		PointerField:    &pointerValue,
		Nested: NestedStruct{
			NestedField: "nested value",
		},
		NilPointer: nil,
	}

	expected := map[string]interface{}{
		"exported_field": "exported value",
		"pointer_field":  "pointer value",
		"nested": map[string]interface{}{
			"nested_field": "nested value",
		},
	}

	result := StructToMap(testStruct)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StructToMap() = %v, want %v", result, expected)
	}
}

func TestStructToMap_NilPointer(t *testing.T) {
	var nilPointer *TestStruct
	result := StructToMap(nilPointer)

	if len(result) != 0 {
		t.Errorf("StructToMap(nil) = %v, want empty map", result)
	}
}

func TestStructToMap_EmptyStruct(t *testing.T) {
	type EmptyStruct struct{}
	result := StructToMap(EmptyStruct{})

	if len(result) != 0 {
		t.Errorf("StructToMap(EmptyStruct{}) = %v, want empty map", result)
	}
}
func TestIsPreservedType(t *testing.T) {
	// Test cases for preserved types
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"Point2D", dbtype.Point2D{}, true},
		{"Pointer to Point2D", &dbtype.Point2D{}, true},
		{"Point3D", dbtype.Point3D{}, true},
		{"Pointer to Point3D", &dbtype.Point3D{}, true},
		{"Time", time.Time{}, true},
		{"LocalDateTime", dbtype.LocalDateTime{}, true},
		{"Date", dbtype.Date{}, true},
		{"Time", dbtype.Time{}, true},
		{"LocalTime", dbtype.LocalTime{}, true},
		{"Duration", dbtype.Duration{}, true},
		{"String", "not preserved", false},
		{"Integer", 123, false},
		{"Struct", struct{}{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPreservedType(tt.input)
			if result != tt.expected {
				t.Errorf("isPreservedType(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
