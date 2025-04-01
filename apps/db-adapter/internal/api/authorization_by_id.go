package api

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
)

type accessMode int

const (
	accessModeNone accessMode = iota
	accessModeRead
	accessModeWrite
)

func userWithIdHasAccessToGivenPerson(ctx context.Context, session neo4j.SessionWithContext, userId, personId int) accessMode {
	resPerson, err := session.ExecuteRead(ctx, memgraph.GetPersonById(ctx, userId))
	if err != nil {
		return accessModeNone
	}

	resPersonMap, ok := resPerson.(map[string]any)
	if !ok {
		return accessModeNone
	}

	if resPersonMap["id"] == userId {
		return accessModeWrite
	}

	AllowAdminAccess, ok := resPersonMap["allow_admin_access"].([]map[string]any)
	if !ok {
		return accessModeNone
	}

	for _, admin := range AllowAdminAccess {
		if admin["id"].(int) == userId {
			return accessModeWrite
		}
	}

	return accessModeNone
}
