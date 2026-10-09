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

//go:embed payloads/update_recipe.json
var update_recipe []byte

func UpdateRecipeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, url, bytes.NewBuffer(update_recipe))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		props, ok := responseBody["Props"].(map[string]any)
		require.True(t, ok, "response should have 'Props'")
		require.Equal(t, "Grandma's Famous Apple Pie", props["name"])
	}
}

func SoftDeleteRecipeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl // test boilerplate
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.Equal(t, "Recipe soft deleted", responseBody["description"])
	}
}

func HardDeleteRecipeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl // test boilerplate
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/hard-delete", dbAdapterUri, recipeId)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.Equal(t, "Recipe hard deleted", responseBody["description"])
	}
}
