package memgraph

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func GetPersonByGoogleId(ctx context.Context, googleId string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetPersonByGoogleIdCypherQuery, map[string]any{
			"google_id": googleId,
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

func UpdatePersonByInviteCode(ctx context.Context, inviteCode string, Person *api.PersonProperties) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, UpdatePersonByInviteCodeCypherQuery, map[string]any{
			"invite_code": inviteCode,
			"props":       *Person,
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
