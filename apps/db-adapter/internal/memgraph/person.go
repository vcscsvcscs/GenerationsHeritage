package memgraph

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func CreatePerson(ctx context.Context, Person *api.PersonProperties) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreatePersonCypherQuery, map[string]any{
			"Person": *Person,
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

func GetPersonById(ctx context.Context, id int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetPersonCypherQuery, map[string]any{
			"id": id,
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
