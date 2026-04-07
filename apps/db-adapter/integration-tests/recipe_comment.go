package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed payloads/recipe_comment.json
var recipeComment []byte

//go:embed payloads/edit_recipe_comment.json
var editRecipeComment []byte

func CommentOnRecipeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/comment", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(recipeComment))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		_, hasComment := responseBody["comment"]
		require.True(t, hasComment, "response should contain 'comment'")

		_, hasCommenter := responseBody["commenter"]
		require.True(t, hasCommenter, "response should contain 'commenter'")
	}
}

func GetRecipeCommentsTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/comment", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		comments, ok := responseBody["comments"].([]any)
		require.True(t, ok, "response should contain 'comments' array")
		require.NotEmpty(t, comments, "should have at least one comment")
	}
}

func UpdateRecipeCommentTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/comment", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, url, bytes.NewBuffer(editRecipeComment))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		comment, ok := responseBody["comment"].(map[string]any)
		require.True(t, ok, "response should contain 'comment' object")

		props, ok := comment["Props"].(map[string]any)
		require.True(t, ok, "comment should have 'Props'")
		require.Equal(t, "This recipe is amazing! Updated with a secret tip.", props["message"])
	}
}

func CommentOnRecipeUpsertTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/comment", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(recipeComment))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify only one comment exists (MERGE should upsert, not duplicate)
		getURL := fmt.Sprintf("%s/recipe/%d/comment", dbAdapterUri, recipeId)
		getReq, err := http.NewRequestWithContext(t.Context(), http.MethodGet, getURL, http.NoBody)
		require.NoError(t, err)
		getReq.Header.Set("X-User-ID", "1")

		getResp, err := client.Do(getReq)
		require.NoError(t, err)
		defer getResp.Body.Close()

		require.Equal(t, http.StatusOK, getResp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(getResp.Body).Decode(&responseBody)
		require.NoError(t, err)

		comments, ok := responseBody["comments"].([]any)
		require.True(t, ok, "response should contain 'comments' array")
		require.Len(t, comments, 1, "MERGE should upsert, not create duplicate comments")
	}
}

func DeleteRecipeCommentTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/comment", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, "Recipe comment deleted", responseBody["description"])
	}
}
