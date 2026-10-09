package integration_tests

import (
	_ "embed"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed payloads/create_recipe_variation.json
var create_recipe_variation []byte

// variationRecipeId is captured from CreateRecipeVariationTest.
var variationRecipeId int //nolint:gochecknoglobals // shared state between sequential integration tests

func CreateRecipeVariationTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/variation", dbAdapterUri, recipeId)

		status, responseBody := doJSON(t, client, http.MethodPost, url, familyMemberId, create_recipe_variation)
		require.Equal(t, http.StatusOK, status)

		recipe, ok := responseBody["recipe"].(map[string]any)
		require.True(t, ok, "response should contain 'recipe' object")
		variationRecipeId = int(recipe["Id"].(float64))
		require.NotEqual(t, recipeId, variationRecipeId)

		variation, ok := responseBody["variation_relationship"].(map[string]any)
		require.True(t, ok, "response should contain 'variation_relationship' object")
		require.Equal(t, "Half the sugar", variation["notes"])
		requireUnixSeconds(t, variation["created_at"])

		_, ok = responseBody["likes_relationship"].(map[string]any)
		require.True(t, ok, "response should contain 'likes_relationship' object")
	}
}

func GetRecipeVariationsTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/variations", dbAdapterUri, recipeId)

		status, responseBody := doJSON(t, client, http.MethodGet, url, 1, nil)
		require.Equal(t, http.StatusOK, status)

		variations, ok := responseBody["variations"].([]any)
		require.True(t, ok, "response should contain 'variations' array")
		require.Len(t, variations, 1)

		entry, ok := variations[0].(map[string]any)
		require.True(t, ok)

		recipe, ok := entry["variation"].(map[string]any)
		require.True(t, ok, "entry should contain the 'variation' recipe")
		require.Equal(t, variationRecipeId, int(recipe["Id"].(float64)))

		relationship, ok := entry["variation_relationship"].(map[string]any)
		require.True(t, ok, "entry should contain 'variation_relationship'")
		require.Equal(t, "Half the sugar", relationship["notes"])
		requireUnixSeconds(t, relationship["created_at"])

		creator, ok := entry["creator"].(map[string]any)
		require.True(t, ok, "entry should contain a flat 'creator'")
		require.InDelta(t, familyMemberId, creator["id"], 0)
		require.NotEmpty(t, creator["first_name"])
	}
}

func GetRecipeVariationsForbiddenTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/variations", dbAdapterUri, recipeId)

		status, _ := doJSON(t, client, http.MethodGet, url, outsiderId, nil)
		require.Equal(t, http.StatusUnauthorized, status)
	}
}
