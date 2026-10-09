package memgraph

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestCreateRecipeVariation(t *testing.T) {
	before := time.Now().Unix()
	var captured map[string]any
	capture := testifymock.MatchedBy(func(params map[string]any) bool {
		captured = params

		return true
	})

	t.Run("success", func(t *testing.T) {
		mockTx := new(mock.Transaction)
		mockTx.On("Run", context.Background(), CreateRecipeVariationCypherQuery, capture).Return(
			recordResult([]string{"recipe", "variation_relationship", "likes_relationship"}, "r", "v", "l"), nil,
		)

		result, err := CreateRecipeVariation(context.Background(), 2, 1, &api.RecipeProperties{}, "less sugar", nil)(mockTx)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"recipe": "r", "variation_relationship": "v", "likes_relationship": "l"}, result)
		require.Equal(t, 2, captured["originalRecipeId"])
		require.Equal(t, 1, captured["creatorId"])
		require.Equal(t, map[string]any{}, captured["LikesProperties"])
		variation, ok := captured["VariationProperties"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, "less sugar", variation["notes"])
		require.GreaterOrEqual(t, variation["created_at"], before, "created_at is a unix timestamp in seconds")
		require.LessOrEqual(t, variation["created_at"], before+60)
	})

	t.Run("single error", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Single", context.Background()).Return(nil, errors.New("single error"))
		mockTx := txRunning(CreateRecipeVariationCypherQuery, testifymock.Anything, mockResult)

		_, err := CreateRecipeVariation(context.Background(), 2, 1, &api.RecipeProperties{}, "", nil)(mockTx)

		require.EqualError(t, err, "single error")
	})
}

func TestGetRecipeVariations(t *testing.T) {
	params := map[string]any{"recipeId": 2}
	keys := []string{"variation", "variation_relationship", "creator"}

	t.Run("collects every variation", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Collect", context.Background()).Return([]*neo4j.Record{
			{Keys: keys, Values: []any{"v1", "r1", "c1"}},
			{Keys: keys, Values: []any{"v2", "r2", "c2"}},
		}, nil)

		result, err := GetRecipeVariations(context.Background(), 2)(txRunning(GetRecipeVariationsCypherQuery, params, mockResult))

		require.NoError(t, err)
		require.Equal(t, map[string]any{"variations": []map[string]any{
			{"variation": "v1", "variation_relationship": "r1", "creator": "c1"},
			{"variation": "v2", "variation_relationship": "r2", "creator": "c2"},
		}}, result)
	})

	t.Run("no variations is an empty list", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Collect", context.Background()).Return([]*neo4j.Record{}, nil)

		result, err := GetRecipeVariations(context.Background(), 2)(txRunning(GetRecipeVariationsCypherQuery, params, mockResult))

		require.NoError(t, err)
		require.Equal(t, map[string]any{"variations": []map[string]any{}}, result)
	})

	t.Run("collect error", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Collect", context.Background()).Return(nil, errors.New("collect error"))

		_, err := GetRecipeVariations(context.Background(), 2)(txRunning(GetRecipeVariationsCypherQuery, params, mockResult))

		require.EqualError(t, err, "collect error")
	})
}
