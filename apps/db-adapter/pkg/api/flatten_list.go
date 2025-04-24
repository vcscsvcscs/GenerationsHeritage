package api

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

func Flatten(input any, uniqueIds *[]int64, result *[]any) error {
	val := reflect.ValueOf(input)
	if uniqueIds == nil {
		uniqueIds = &[]int64{}
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := range val.Len() {
			if err := Flatten(val.Index(i).Interface(), uniqueIds, result); err != nil {
				return err
			}
		}
	case reflect.Map:
		id, ok := input.(map[string]any)["id"].(int64)
		if !ok {
			return fmt.Errorf("could not convert id to int: %v", input.(map[string]any)["id"])
		}

		if !slices.Contains(*uniqueIds, id) {
			*result = append(*result, input.(map[string]any))
			*uniqueIds = append(*uniqueIds, id)
		}
	case reflect.Struct:
		switch input.(type) {
		case dbtype.Node:
			node := val.Interface().(dbtype.Node)
			if !slices.Contains(*uniqueIds, node.Id) { //nolint:staticcheck // this is a known issue with the neo4j-go-driver
				*result = append(*result, node)
				*uniqueIds = append(*uniqueIds, node.Id) //nolint:staticcheck // this is a known issue with the neo4j-go-driver
			}
		case dbtype.Relationship:
			relationship := val.Interface().(dbtype.Relationship)
			if !slices.Contains(*uniqueIds, relationship.Id) { //nolint:staticcheck // this is a known issue with the neo4j-go-driver
				*result = append(*result, relationship)
				*uniqueIds = append(*uniqueIds, relationship.Id) //nolint:staticcheck // this is a known issue with the neo4j-go-driver
			}
		}
	default:
		return fmt.Errorf("unexpected type: %T", input)
	}

	return nil
}
