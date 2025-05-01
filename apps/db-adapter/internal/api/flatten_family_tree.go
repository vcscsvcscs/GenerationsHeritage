package api

import (
	"fmt"

	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

type FamilyTree struct {
	People        []any `json:"people"`
	Relationships []any `json:"relationships"`
}

func FlattenFamilyTree(input any, result *FamilyTree) error {
	root, ok := input.(map[string]any)
	if !ok {
		return fmt.Errorf("could not convert result to map[string]any")
	}

	var uniqueIds []int64
	err := api.Flatten(root["people"], &uniqueIds, &result.People)
	if err != nil {
		return err
	}

	uniqueIds = []int64{}
	return api.Flatten(root["relationships"], &uniqueIds, &result.Relationships)
}
