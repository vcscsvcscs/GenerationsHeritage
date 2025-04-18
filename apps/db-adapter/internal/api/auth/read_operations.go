package auth

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func CouldSeePersonsProfile(ctx context.Context, session neo4j.SessionWithContext, userId, XUserID int) error {
	if CouldManagePersonUnknownAdmin(ctx, session, userId, XUserID) == nil {
		return nil
	}

	res, err := session.ExecuteRead(ctx, memgraph.GetFamilyTreeById(ctx, XUserID))
	if err != nil {
		return err
	}

	resMap, ok := res.(map[string]any)
	if !ok {
		return fmt.Errorf("could not convert result to map[string]any")
	}

	people, ok := resMap["people"].([]api.OptimizedPersonNode)
	if !ok {
		return fmt.Errorf("could not convert people to []api.PersonProperties")
	}

	for _, person := range people {
		if *person.Id == userId {
			return nil
		}
	}

	return fmt.Errorf("user %d does not have permission to see user %d", XUserID, userId)
}
