package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed payloads/create_recipe.json
var create_recipe []byte

// recipeId is captured from CreateRecipeForPersonTest and used by subsequent recipe tests.
// It is set dynamically because memgraph assigns node IDs sequentially.
var recipeId int //nolint:gochecknoglobals // shared state between sequential integration tests

func CreateRecipeForPersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/1/recipes"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(create_recipe))
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

		recipe, ok := responseBody["recipe"].(map[string]any)
		require.True(t, ok, "response should contain 'recipe' object")

		id, ok := recipe["Id"]
		require.True(t, ok, "recipe should have an 'Id' field")
		recipeId = int(id.(float64))

		props, ok := recipe["Props"].(map[string]any)
		require.True(t, ok, "recipe should have 'Props'")
		require.Equal(t, "Grandma's Apple Pie", props["name"])
		require.Equal(t, "Dessert", props["category"])

		_, hasRel := responseBody["relationship"]
		require.True(t, hasRel, "response should contain 'relationship' object")
	}
}

func GetRecipesByPersonIdTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		status, responseBody := doJSON(t, client, http.MethodGet, dbAdapterUri+"/person/1/recipes", 1, nil)
		require.Equal(t, http.StatusOK, status)

		entries, ok := responseBody["entries"].([]any)
		require.True(t, ok, "response should contain 'entries' array")
		require.NotEmpty(t, entries, "person should have at least one recipe")

		entry := entryForRecipe(entries, recipeId)
		require.NotNil(t, entry, "the created recipe should be listed")
		require.NotNil(t, entry["relationship"], "the creator also likes the recipe")
		require.Equal(t, true, entry["created"])
		require.Equal(t, true, entry["can_edit"])
	}
}

func GetFamilyCookbookTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		status, responseBody := doJSON(t, client, http.MethodGet, dbAdapterUri+"/cookbook?distance=3", 1, nil)
		require.Equal(t, http.StatusOK, status)

		entries, ok := responseBody["entries"].([]any)
		require.True(t, ok, "response should contain 'entries' array")
		require.NotEmpty(t, entries, "cookbook should have at least one entry")

		seen := map[int]bool{}
		for _, raw := range entries {
			entry, ok := raw.(map[string]any)
			require.True(t, ok, "entry should be an object")

			recipe, ok := entry["recipe"].(map[string]any)
			require.True(t, ok, "entry should have a non-null recipe")
			id := int(recipe["Id"].(float64))
			require.False(t, seen[id], "recipe %d should be listed once", id)
			seen[id] = true

			_, ok = entry["added_by"].(map[string]any)
			require.True(t, ok, "entry should have a non-null added_by")
			_, ok = entry["can_edit"].(bool)
			require.True(t, ok, "entry should say whether the recipe can be edited")
		}

		own := entryForRecipe(entries, recipeId)
		require.NotNil(t, own, "the user's own recipe should be in the cookbook")
		require.InDelta(t, 1, own["added_by"].(map[string]any)["id"], 0, "the user is preferred as added_by")
		require.Equal(t, true, own["can_edit"])

		variation := entryForRecipe(entries, variationRecipeId)
		require.NotNil(t, variation, "a family member's recipe should be in the cookbook")
		require.InDelta(t, familyMemberId, variation["added_by"].(map[string]any)["id"], 0)
		require.Equal(t, false, variation["can_edit"])
	}
}
