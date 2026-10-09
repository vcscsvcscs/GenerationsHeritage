package memgraph

import (
	"context"
	"errors"
	"testing"
	"time"

	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
)

func TestCommentOnRecipe(t *testing.T) {
	before := time.Now().Unix()
	var captured map[string]any
	capture := testifymock.MatchedBy(func(params map[string]any) bool {
		captured = params

		return true
	})
	mockTx := new(mock.Transaction)
	mockTx.On("Run", context.Background(), CommentOnRecipeCypherQuery, capture).
		Return(recordResult([]string{"comment", "commenter"}, "commentValue", "commenterValue"), nil)

	result, err := CommentOnRecipe(context.Background(), 1, 2, "yummy")(mockTx)

	require.NoError(t, err)
	require.Equal(t, map[string]any{"comment": "commentValue", "commenter": "commenterValue"}, result)
	require.Equal(t, 1, captured["personId"])
	require.Equal(t, 2, captured["recipeId"])
	comment, ok := captured["Comment"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "yummy", comment["message"])
	require.GreaterOrEqual(t, comment["sent_at"], before, "sent_at is a unix timestamp in seconds")
	require.LessOrEqual(t, comment["sent_at"], before+60)
}

func TestGetRecipeComments(t *testing.T) {
	params := map[string]any{"recipeId": 2}

	t.Run("success", func(t *testing.T) {
		mockTx := txRunning(GetRecipeCommentsCypherQuery, params, recordResult([]string{"comments"}, []any{"c1"}))

		result, err := GetRecipeComments(context.Background(), 2)(mockTx)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"comments": []any{"c1"}}, result)
	})

	t.Run("single error", func(t *testing.T) {
		mockResult := new(mock.Result)
		mockResult.On("Single", context.Background()).Return(nil, errors.New("single error"))

		_, err := GetRecipeComments(context.Background(), 2)(txRunning(GetRecipeCommentsCypherQuery, params, mockResult))

		require.EqualError(t, err, "single error")
	})
}

func TestUpdateRecipeComment(t *testing.T) {
	before := time.Now().Unix()
	var captured map[string]any
	capture := testifymock.MatchedBy(func(params map[string]any) bool {
		captured = params

		return true
	})

	t.Run("success", func(t *testing.T) {
		mockTx := new(mock.Transaction)
		mockTx.On("Run", context.Background(), UpdateRecipeCommentCypherQuery, capture).
			Return(recordResult([]string{"comment", "commenter"}, "commentValue", "commenterValue"), nil)

		result, err := UpdateRecipeComment(context.Background(), 1, 2, "edited")(mockTx)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"comment": "commentValue", "commenter": "commenterValue"}, result)
		require.Equal(t, 1, captured["personId"])
		require.Equal(t, 2, captured["recipeId"])
		require.Equal(t, "edited", captured["message"])
		require.GreaterOrEqual(t, captured["edited"], before, "edited is a unix timestamp in seconds")
		require.LessOrEqual(t, captured["edited"], before+60)
	})

	t.Run("missing comment is not found", func(t *testing.T) {
		mockTx := new(mock.Transaction)
		mockTx.On("Run", context.Background(), UpdateRecipeCommentCypherQuery, testifymock.Anything).Return(emptyResult(), nil)

		result, err := UpdateRecipeComment(context.Background(), 1, 2, "edited")(mockTx)

		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, result)
	})

	t.Run("run error", func(t *testing.T) {
		mockTx := new(mock.Transaction)
		mockTx.On("Run", context.Background(), UpdateRecipeCommentCypherQuery, testifymock.Anything).Return(nil, errors.New("run error"))

		_, err := UpdateRecipeComment(context.Background(), 1, 2, "edited")(mockTx)

		require.EqualError(t, err, "run error")
	})
}

func TestDeleteRecipeComment(t *testing.T) {
	params := map[string]any{"personId": 1, "recipeId": 2}

	t.Run("success", func(t *testing.T) {
		mockTx := txRunning(DeleteRecipeCommentCypherQuery, params, recordResult([]string{"deleted"}, true))

		result, err := DeleteRecipeComment(context.Background(), 1, 2)(mockTx)

		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("missing comment is not found", func(t *testing.T) {
		_, err := DeleteRecipeComment(context.Background(), 1, 2)(txRunning(DeleteRecipeCommentCypherQuery, params, emptyResult()))

		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("run error", func(t *testing.T) {
		mockTx := new(mock.Transaction)
		mockTx.On("Run", context.Background(), DeleteRecipeCommentCypherQuery, params).Return(nil, errors.New("run error"))

		_, err := DeleteRecipeComment(context.Background(), 1, 2)(mockTx)

		require.EqualError(t, err, "run error")
	})
}
