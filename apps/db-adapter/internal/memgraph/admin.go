package memgraph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateAdminRelationship(ctx context.Context, userId int, adminId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateAdminRelationshipCypherQuery, map[string]any{
			"id1": adminId,
			"id2": userId,
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

func DeleteAdminRelationship(ctx context.Context, userId int, adminId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, DeleteAdminRelationshipCypherQuery, map[string]any{
			"id1": adminId,
			"id2": userId,
		})
		if err != nil {
			return nil, err
		}

		if result.Peek(ctx) {
			return nil, fmt.Errorf("there was a returned value, when deleting admin but there should be none")
		}

		return nil, nil
	}
}

func GetAdminRelationship(ctx context.Context, userId int, adminId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetAdminRelationshipCypherQuery, map[string]any{
			"id1": adminId,
			"id2": userId,
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

func GetProfileAdmins(ctx context.Context, userId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetProfileAdminsCypherQuery, map[string]any{
			"id": userId,
		})
		if err != nil {
			return nil, err
		}

		return result.Collect(ctx)
	}
}
