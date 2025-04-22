package auth

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func CouldSeePersonsProfile(ctx context.Context, session neo4j.SessionWithContext, userId, xUserID int) error {
	if CouldManagePersonUnknownAdmin(ctx, session, userId, xUserID) == nil {
		return nil
	}

	res, err := session.ExecuteRead(ctx, memgraph.GetFamilyTreeById(ctx, xUserID))
	if err != nil {
		return err
	}

	resMap, ok := res.(map[string]any)
	if !ok {
		return fmt.Errorf("could not convert result to map[string]any")
	}

	var uniqueIds []int64
	var flattenedPeople []any
	if err := api.Flatten(resMap["people"], &uniqueIds, &flattenedPeople); err != nil {
		return fmt.Errorf("could not convert people to []map[string]any: %w", err)
	}

	for _, person := range flattenedPeople {
		person, ok := person.(map[string]any)
		if !ok {
			return fmt.Errorf("could not convert person to map[string]any")
		}

		if person["id"].(int64) == int64(userId) {
			return nil
		}
	}

	return fmt.Errorf("user %d does not have permission to see user %d", xUserID, userId)
}
