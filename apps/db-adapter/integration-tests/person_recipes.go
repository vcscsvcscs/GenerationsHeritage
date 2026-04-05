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
		url := dbAdapterUri + "/person/1/recipes"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		recipes, ok := responseBody["recipes"].([]any)
		require.True(t, ok, "response should contain 'recipes' array")
		require.NotEmpty(t, recipes, "person should have at least one recipe")

		recipeRelations, ok := responseBody["recipeRelations"].([]any)
		require.True(t, ok, "response should contain 'recipeRelations' array")
		require.NotEmpty(t, recipeRelations, "person should have at least one recipe relation")
	}
}

func GetFamilyCookbookTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/cookbook?distance=3"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("X-User-ID", "1")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		entries, ok := responseBody["entries"].([]any)
		require.True(t, ok, "response should contain 'entries' array")
		require.NotEmpty(t, entries, "cookbook should have at least one entry")
	}
}
