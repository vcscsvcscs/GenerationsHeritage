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
	ExportedField   string `json:"exported_field"`
	IgnoredField    string `json:"-"`
	unexportedField string
	PointerField    *string      `json:"pointer_field"`
	NilPointer      *string      `json:"nil_pointer"`
	Nested          NestedStruct `json:"nested"`
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

	expected := map[string]any{
		"exported_field": "exported value",
		"pointer_field":  "pointer value",
		"nested": map[string]any{
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
		input    any
		name     string
		expected bool
	}{
		{dbtype.Point2D{}, "Point2D", true},
		{&dbtype.Point2D{}, "Pointer to Point2D", true},
		{dbtype.Point3D{}, "Point3D", true},
		{&dbtype.Point3D{}, "Pointer to Point3D", true},
		{time.Time{}, "Time", true},
		{dbtype.LocalDateTime{}, "LocalDateTime", true},
		{dbtype.Date{}, "Date", true},
		{dbtype.Time{}, "Time", true},
		{dbtype.LocalTime{}, "LocalTime", true},
		{dbtype.Duration{}, "Duration", true},
		{"not preserved", "String", false},
		{123, "Integer", false},
		{struct{}{}, "Struct", false},
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
