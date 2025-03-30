package memgraph

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetRelationship(ctx context.Context, id1, id2 int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetRelationshipCypherQuery, map[string]any{
			"id1": id1,
			"id2": id2,
		})
		if err != nil {
			return nil, err
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		return records, nil
	}
}

func CreateChildParentRelationship(ctx context.Context, childId, parentId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateChildParentRelationshipCypherQuery, map[string]any{
			"childId":  childId,
			"parentId": parentId,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return record.AsMap(), nil
	}
}
