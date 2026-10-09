package integration_tests

import (
	_ "embed"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed payloads/create_recipe_relationship.json
var create_recipe_relationship []byte

//go:embed payloads/create_recipe_relationship_for_person.json
var create_recipe_relationship_for_person []byte

// CreateRecipeRelationshipTest has a family member like the recipe; the person defaults to the caller.
func CreateRecipeRelationshipTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/relationship", dbAdapterUri, recipeId)

		status, responseBody := doJSON(t, client, http.MethodPost, url, familyMemberId, create_recipe_relationship)
		require.Equal(t, http.StatusOK, status)

		_, ok := responseBody["Id"]
		require.True(t, ok, "response should contain relationship 'Id'")
		require.InDelta(t, familyMemberId, responseBody["StartId"], 0, "the caller is the one liking")
	}
}

func CreateRecipeRelationshipForbiddenTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl // test boilerplate
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/relationship", dbAdapterUri, recipeId)

		t.Run("cannot like on behalf of a person one does not manage", func(t *testing.T) {
			status, _ := doJSON(t, client, http.MethodPost, url, familyMemberId, create_recipe_relationship_for_person)
			require.Equal(t, http.StatusUnauthorized, status)
		})

		t.Run("cannot like a recipe one cannot see", func(t *testing.T) {
			status, _ := doJSON(t, client, http.MethodPost, url, outsiderId, create_recipe_relationship)
			require.Equal(t, http.StatusUnauthorized, status)
		})
	}
}

func DeleteRecipeRelationshipForbiddenTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("cannot remove another person's like", func(t *testing.T) {
			url := fmt.Sprintf("%s/recipe/%d/relationship?personId=%d", dbAdapterUri, recipeId, familyMemberId)

			status, _ := doJSON(t, client, http.MethodDelete, url, 1, nil)
			require.Equal(t, http.StatusUnauthorized, status)
		})

		t.Run("cannot remove a like from a recipe one cannot see", func(t *testing.T) {
			url := fmt.Sprintf("%s/recipe/%d/relationship?personId=%d", dbAdapterUri, recipeId, outsiderId)

			status, _ := doJSON(t, client, http.MethodDelete, url, outsiderId, nil)
			require.Equal(t, http.StatusUnauthorized, status)
		})
	}
}

func DeleteRecipeRelationshipTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/relationship?personId=%d", dbAdapterUri, recipeId, familyMemberId)

		status, responseBody := doJSON(t, client, http.MethodDelete, url, familyMemberId, nil)
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, "Recipe relationship deleted", responseBody["description"])
	}
}

// CreatorUnlikingKeepsRecipeTest checks that unliking a self created recipe neither orphans nor hides it.
func CreatorUnlikingKeepsRecipeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/relationship?personId=1", dbAdapterUri, recipeId)
		status, _ := doJSON(t, client, http.MethodDelete, url, 1, nil)
		require.Equal(t, http.StatusOK, status)

		status, responseBody := doJSON(t, client, http.MethodGet, dbAdapterUri+"/person/1/recipes", 1, nil)
		require.Equal(t, http.StatusOK, status)

		entries, ok := responseBody["entries"].([]any)
		require.True(t, ok)
		entry := entryForRecipe(entries, recipeId)
		require.NotNil(t, entry, "the creator should still see the recipe they created")
		require.Nil(t, entry["relationship"], "the creator no longer likes it")
		require.Equal(t, true, entry["created"])
		require.Equal(t, true, entry["can_edit"], "the creator can still edit it")

		status, responseBody = doJSON(t, client, http.MethodGet, dbAdapterUri+"/cookbook?distance=3", 1, nil)
		require.Equal(t, http.StatusOK, status)

		entries, ok = responseBody["entries"].([]any)
		require.True(t, ok)
		require.NotNil(t, entryForRecipe(entries, recipeId), "the recipe stays in the cookbook")
	}
}
