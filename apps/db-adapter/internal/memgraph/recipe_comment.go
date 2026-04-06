package memgraph

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type RecipeComment struct {
	SentAt  time.Time  `json:"sent_at"`
	Edited  *time.Time `json:"edited,omitempty"`
	Message string     `json:"message"`
}

func CommentOnRecipe(
	ctx context.Context,
	personId, recipeId int,
	message string,
) neo4j.ManagedTransactionWork {
	comment := map[string]any{
		"message": message,
		"sent_at": time.Now().Unix(),
	}

	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CommentOnRecipeCypherQuery, map[string]any{
			"personId": personId,
			"recipeId": recipeId,
			"Comment":  comment,
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

func GetRecipeComments(ctx context.Context, recipeId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetRecipeCommentsCypherQuery, map[string]any{
			"recipeId": recipeId,
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

func UpdateRecipeComment(
	ctx context.Context,
	personId, recipeId int,
	message string,
) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, UpdateRecipeCommentCypherQuery, map[string]any{
			"personId": personId,
			"recipeId": recipeId,
			"message":  message,
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

func DeleteRecipeComment(ctx context.Context, personId, recipeId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, DeleteRecipeCommentCypherQuery, map[string]any{
			"personId": personId,
			"recipeId": recipeId,
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
