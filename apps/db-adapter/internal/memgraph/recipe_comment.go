package memgraph

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

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
	edited := time.Now().Unix()

	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, UpdateRecipeCommentCypherQuery, map[string]any{
			"personId": personId,
			"recipeId": recipeId,
			"message":  message,
			"edited":   edited,
		})
		if err != nil {
			return nil, err
		}

		record, err := singleOrNotFound(ctx, result)
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

		if _, err := singleOrNotFound(ctx, result); err != nil {
			return nil, err
		}

		return nil, nil
	}
}
