package memgraph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func CreatePerson(ctx context.Context, person *api.PersonProperties) neo4j.ManagedTransactionWork {
	convertedPerson := StructToMap(person)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreatePersonCypherQuery, map[string]any{
			"Person": convertedPerson,
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

		person, ok := record.Get("person")
		if !ok {
			return nil, fmt.Errorf("person not found")
		}

		return person, nil
	}
}

func UpdatePerson(ctx context.Context, id int, person *api.PersonProperties) neo4j.ManagedTransactionWork {
	convertedPerson := StructToMap(person)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, UpdatePersonCypherQuery, map[string]any{
			"id":    id,
			"props": convertedPerson,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		person, ok := record.Get("person")
		if !ok {
			return nil, fmt.Errorf("person not found")
		}

		return person, nil
	}
}

func SoftDeletePerson(ctx context.Context, id int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, SoftDeletePersonCypherQuery, map[string]any{
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

func HardDeletePerson(ctx context.Context, id int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, HardDeletePersonCypherQuery, map[string]any{
			"id": id,
		})
		if err != nil {
			return nil, err
		}

		if result.Peek(ctx) {
			return nil, fmt.Errorf("record was returned when it wasn't supposed to happen")
		}

		return nil, nil
	}
}
