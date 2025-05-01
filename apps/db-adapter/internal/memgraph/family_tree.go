package memgraph

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// returned map has "people" which is a slice of OptimizedPersonNode and relationships which a slice of Relatioship type.
func GetFamilyTreeById(ctx context.Context, userId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetBloodRelativesCypherQuery, map[string]any{
			"id": userId,
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

// returned map has "people" which is a slice of OptimizedPersonNode and relationships which a slice of Relatioship type.
func GetFamilyTreeWithSpousesById(ctx context.Context, userId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetFamilyTreeWithSpousesCypherQuery, map[string]any{
			"id": userId,
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
