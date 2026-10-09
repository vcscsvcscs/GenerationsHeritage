package memgraph

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func recordResult(keys []string, values ...any) *mock.Result {
	mockResult := new(mock.Result)
	mockResult.On("Peek", context.Background()).Return(true)
	mockResult.On("Single", context.Background()).Return(&neo4j.Record{Keys: keys, Values: values}, nil)

	return mockResult
}

func emptyResult() *mock.Result {
	mockResult := new(mock.Result)
	mockResult.On("Peek", context.Background()).Return(false)
	mockResult.On("Err").Return(nil)

	return mockResult
}

func txRunning(query string, params any, result *mock.Result) *mock.Transaction {
	mockTx := new(mock.Transaction)
	mockTx.On("Run", context.Background(), query, params).Return(result, nil)

	return mockTx
}

func TestSingleOrNotFound(t *testing.T) {
	t.Run("returns the record", func(t *testing.T) {
		record, err := singleOrNotFound(context.Background(), recordResult([]string{"k"}, "v"))

		require.NoError(t, err)
		require.Equal(t, map[string]any{"k": "v"}, record.AsMap())
	})

	t.Run("empty result is not found", func(t *testing.T) {
		_, err := singleOrNotFound(context.Background(), emptyResult())

		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("result error is returned", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Peek", context.Background()).Return(false)
		mockResult.On("Err").Return(errors.New("stream error"))

		_, err := singleOrNotFound(context.Background(), mockResult)

		require.EqualError(t, err, "stream error")
		require.NotErrorIs(t, err, ErrNotFound)
	})
}

func TestGetRecipe(t *testing.T) {
	params := map[string]any{"recipeId": 2, "userId": 1}

	t.Run("success", func(t *testing.T) {
		mockTx := txRunning(GetRecipeByIdCypherQuery, params, recordResult([]string{"recipe", "can_edit"}, "recipeValue", true))

		result, err := GetRecipe(context.Background(), 2, 1)(mockTx)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"recipe": "recipeValue", "can_edit": true}, result)
	})

	t.Run("not found", func(t *testing.T) {
		result, err := GetRecipe(context.Background(), 2, 1)(txRunning(GetRecipeByIdCypherQuery, params, emptyResult()))

		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, result)
	})

	t.Run("run error", func(t *testing.T) {
		mockTx := new(mock.Transaction)
		mockTx.On("Run", context.Background(), GetRecipeByIdCypherQuery, params).Return(nil, errors.New("run error"))

		result, err := GetRecipe(context.Background(), 2, 1)(mockTx)

		require.EqualError(t, err, "run error")
		require.Nil(t, result)
	})
}

func TestUpdateAndSoftDeleteRecipeNotFound(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		mockTx := txRunning(UpdateRecipeCypherQuery, map[string]any{"id": 123, "props": map[string]any{}}, emptyResult())

		_, err := UpdateRecipe(context.Background(), 123, &api.RecipeProperties{})(mockTx)

		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("soft delete of an already deleted recipe", func(t *testing.T) {
		mockTx := txRunning(SoftDeleteRecipeCypherQuery, map[string]any{"id": 123}, emptyResult())

		_, err := SoftDeleteRecipe(context.Background(), 123)(mockTx)

		require.ErrorIs(t, err, ErrNotFound)
	})
}

func TestGetFamilyCookbook(t *testing.T) {
	query := fmt.Sprintf(GetFamilyCookbookCypherQueryTemplate, 3)
	params := map[string]any{"id": 7}

	t.Run("embeds the distance in the query", func(t *testing.T) {
		require.Contains(t, query, "*1..3")
		require.NotContains(t, query, "%d")
	})

	t.Run("success", func(t *testing.T) {
		mockTx := txRunning(query, params, recordResult([]string{"entries"}, []any{"entry"}))

		result, err := GetFamilyCookbook(context.Background(), 7, 3)(mockTx)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"entries": []any{"entry"}}, result)
	})

	t.Run("single error", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Single", context.Background()).Return(nil, errors.New("single error"))

		result, err := GetFamilyCookbook(context.Background(), 7, 3)(txRunning(query, params, mockResult))

		require.EqualError(t, err, "single error")
		require.Nil(t, result)
	})
}

func TestCouldManageRecipe(t *testing.T) {
	params := map[string]any{"recipeId": 2, "userId": 1}

	t.Run("allowed", func(t *testing.T) {
		mockTx := txRunning(CouldManageRecipeCypherQuery, params, recordResult([]string{"r"}, "recipeValue"))

		result, err := CouldManageRecipe(context.Background(), 2, 1)(mockTx)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"r": "recipeValue"}, result)
	})

	t.Run("denied", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Single", context.Background()).Return(nil, errors.New("no records"))

		result, err := CouldManageRecipe(context.Background(), 2, 1)(txRunning(CouldManageRecipeCypherQuery, params, mockResult))

		require.EqualError(t, err, "user 1 does not have permission to manage recipe 2")
		require.Nil(t, result)
	})
}

func TestCouldSeeRecipe(t *testing.T) {
	query := fmt.Sprintf(CouldSeeRecipeCypherQueryTemplate, MaxCookbookDistance)
	params := map[string]any{"recipeId": 2, "userId": 1}

	t.Run("reaches as far as the cookbook does", func(t *testing.T) {
		require.Contains(t, query, fmt.Sprintf("*1..%d", MaxCookbookDistance))
	})

	t.Run("allowed", func(t *testing.T) {
		mockTx := txRunning(query, params, recordResult([]string{"r"}, "recipeValue"))

		result, err := CouldSeeRecipe(context.Background(), 2, 1)(mockTx)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"r": "recipeValue"}, result)
	})

	t.Run("denied", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Single", context.Background()).Return(nil, errors.New("no records"))

		result, err := CouldSeeRecipe(context.Background(), 2, 1)(txRunning(query, params, mockResult))

		require.EqualError(t, err, "user 1 does not have permission to see recipe 2")
		require.Nil(t, result)
	})
}
