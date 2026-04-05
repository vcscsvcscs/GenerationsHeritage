package memgraph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func CreateRecipeForPerson(
	ctx context.Context, personId int, recipe *api.RecipeProperties, relationship *api.LikesProperties,
) neo4j.ManagedTransactionWork {
	convertedRecipe := StructToMap(recipe)
	convertedRelationship := map[string]any{}
	if relationship != nil {
		convertedRelationship = StructToMap(relationship)
	}
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateRecipeWithRelationshipCypherQuery, map[string]any{
			"personId":     personId,
			"Recipe":       convertedRecipe,
			"Relationship": convertedRelationship,
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

func GetRecipesByPersonId(ctx context.Context, personId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, GetRecipesByPersonIdCypherQuery, map[string]any{
			"id": personId,
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

func UpdateRecipe(ctx context.Context, id int, recipe *api.RecipeProperties) neo4j.ManagedTransactionWork {
	convertedRecipe := StructToMap(recipe)
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, UpdateRecipeCypherQuery, map[string]any{
			"id":    id,
			"props": convertedRecipe,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		recipe, ok := record.Get("recipe")
		if !ok {
			return nil, fmt.Errorf("recipe not found")
		}

		return recipe, nil
	}
}

func SoftDeleteRecipe(ctx context.Context, id int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, SoftDeleteRecipeCypherQuery, map[string]any{
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

func HardDeleteRecipe(ctx context.Context, id int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, HardDeleteRecipeCypherQuery, map[string]any{
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

func CreateRecipeRelationship(
	ctx context.Context, personId, recipeId int, relationship *api.LikesProperties,
) neo4j.ManagedTransactionWork {
	convertedRelationship := map[string]any{}
	if relationship != nil {
		convertedRelationship = StructToMap(relationship)
	}
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CreateRecipeRelationshipCypherQuery, map[string]any{
			"personId":     personId,
			"recipeId":     recipeId,
			"Relationship": convertedRelationship,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		rel, ok := record.Get("relationship")
		if !ok {
			return nil, fmt.Errorf("relationship not found")
		}

		return rel, nil
	}
}

func DeleteRecipeRelationship(ctx context.Context, personId, recipeId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, DeleteRecipeRelationshipCypherQuery, map[string]any{
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

// CouldManageRecipe checks whether a user has permission to manage (update/delete) a recipe.
// A user can manage a recipe if they have a Created relationship to it, or if they are
// an admin of a person who has a Created relationship to it.
func CouldManageRecipe(ctx context.Context, recipeId, userId int) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, CouldManageRecipeCypherQuery, map[string]any{
			"recipeId": recipeId,
			"userId":   userId,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, fmt.Errorf("user %d does not have permission to manage recipe %d", userId, recipeId)
		}

		return record.AsMap(), nil
	}
}

// GetFamilyCookbook retrieves all recipes from family members within a given distance.
// Distance is the number of relationship hops to traverse.
func GetFamilyCookbook(ctx context.Context, userId, distance int) neo4j.ManagedTransactionWork {
	// Cypher does not support parameterized variable-length relationship bounds,
	// so we build the query with the distance embedded safely as an integer.
	query := fmt.Sprintf(GetFamilyCookbookCypherQueryTemplate, distance)

	return func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, map[string]any{
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
