package memgraph

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
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

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return record.AsMap(), nil
	}
}

func DeleteRelationship(ctx context.Context, id1, id2 int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, DeleteRelationshipCypherQuery, map[string]any{
			"id1": id1,
			"id2": id2,
		})
		if err != nil {
			return nil, err
		}

		return !result.Peek(ctx), nil
	}
}

func UpdateRelationship(
	ctx context.Context, id1, id2 int, relationship api.FamilyRelationship,
) neo4j.ManagedTransactionWork {
	convertedRelationship := StructToMap(relationship)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, UpdateRelationshipCypherQuery, map[string]any{
			"id1":          id1,
			"id2":          id2,
			"relationship": convertedRelationship,
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

func CreateChildParentRelationship(
	ctx context.Context, childId, parentId int, relationship api.FamilyRelationship,
) neo4j.ManagedTransactionWork {
	convertedRelationship := StructToMap(relationship)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateChildParentRelationshipCypherQuery, map[string]any{
			"childId":            childId,
			"parentId":           parentId,
			"childRelationship":  convertedRelationship,
			"parentRelationship": convertedRelationship,
		})
		if err != nil {
			return nil, err
		}

		return result.Single(ctx)
	}
}

func CreateSiblingRelationship(
	ctx context.Context, siblingId1, siblingId2 int, relationship api.FamilyRelationship,
) neo4j.ManagedTransactionWork {
	convertedRelationship := StructToMap(relationship)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateSiblingRelationshipCypherQuery, map[string]any{
			"id1":           siblingId1,
			"id2":           siblingId2,
			"Relationship1": convertedRelationship,
			"Relationship2": convertedRelationship,
		})
		if err != nil {
			return nil, err
		}

		return result.Collect(ctx)
	}
}

func CreateSpouseRelationship(
	ctx context.Context, spouseId1, spouseId2 int, relationship api.FamilyRelationship,
) neo4j.ManagedTransactionWork {
	convertedRelationship := StructToMap(relationship)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateSpouseRelationshipCypherQuery, map[string]any{
			"id1":           spouseId1,
			"id2":           spouseId2,
			"Relationship1": convertedRelationship,
			"Relationship2": convertedRelationship,
		})
		if err != nil {
			return nil, err
		}

		return result.Collect(ctx)
	}
}
