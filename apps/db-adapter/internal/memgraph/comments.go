package memgraph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func UpsertCommentOnProfile(ctx context.Context, commenter, profile int, comment *api.Message) neo4j.ManagedTransactionWork {
	convertedComment := StructToMap(comment)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CommentCypherQuery, map[string]any{
			"id1":     commenter,
			"id2":     profile,
			"comment": convertedComment,
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

func GetCommentsOnProfile(ctx context.Context, profile int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CommentsOnProfileCypherQuery, map[string]any{
			"id": profile,
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

func DeleteComment(ctx context.Context, commenter, profile int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, DeleteCommentCypherQuery, map[string]any{
			"id1": commenter,
			"id2": profile,
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
