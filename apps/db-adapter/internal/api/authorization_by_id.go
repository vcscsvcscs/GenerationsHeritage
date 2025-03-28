package api

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
)

func userWithIdHasAccessToGivenPerson(ctx context.Context, session neo4j.SessionWithContext, userId, personId int) bool {
	_, err := session.ExecuteRead(ctx, memgraph.GetPersonById(ctx, userId))
	//TODO
	return err == nil
}
