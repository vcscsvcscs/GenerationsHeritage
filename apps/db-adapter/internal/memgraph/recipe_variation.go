package memgraph

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func CreateRecipeVariation(
	ctx context.Context,
	originalRecipeId, creatorId int,
	recipeProps *api.RecipeProperties,
	variationNotes string,
	likesProps *api.LikesProperties,
) neo4j.ManagedTransactionWork {
	convertedRecipe := StructToMap(recipeProps)
	variationProperties := map[string]any{
		"notes":      variationNotes,
		"created_at": time.Now().UnixMilli(),
	}
	
	var convertedLikes map[string]any
	if likesProps != nil {
		convertedLikes = StructToMap(likesProps)
	}

	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateRecipeVariationCypherQuery, map[string]any{
			"originalRecipeId":    originalRecipeId,
			"creatorId":           creatorId,
			"RecipeProperties":    convertedRecipe,
			"VariationProperties": variationProperties,
			"LikesProperties":     convertedLikes,
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

func GetRecipeVariations(ctx context.Context, recipeId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetRecipeVariationsCypherQuery, map[string]any{
			"recipeId": recipeId,
		})
		if err != nil {
			return nil, err
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		variations := make([]map[string]any, 0, len(records))
		for _, record := range records {
			variations = append(variations, record.AsMap())
		}

		return map[string]any{
			"variations": variations,
		}, nil
	}
}

func CouldSeeRecipe(ctx context.Context, recipeId, userId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CouldSeeRecipeCypherQuery, map[string]any{
			"recipeId": recipeId,
			"userId":   userId,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, fmt.Errorf("user %d does not have permission to see recipe %d", userId, recipeId)
		}

		return record.AsMap(), nil
	}
}
