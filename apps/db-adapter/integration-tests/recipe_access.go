package integration_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRecipeForPersonForbiddenTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		// a family member who is not an admin of person 1 must not create recipes on their behalf
		status, _ := doJSON(t, client, http.MethodPost, dbAdapterUri+"/person/1/recipes", familyMemberId, create_recipe)
		require.Equal(t, http.StatusUnauthorized, status)
	}
}

func GetRecipeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d", dbAdapterUri, recipeId)

		t.Run("creator can edit", func(t *testing.T) {
			status, responseBody := doJSON(t, client, http.MethodGet, url, 1, nil)
			require.Equal(t, http.StatusOK, status)

			recipe, ok := responseBody["recipe"].(map[string]any)
			require.True(t, ok, "response should contain 'recipe' object")
			require.InDelta(t, recipeId, recipe["Id"], 0)
			require.Equal(t, true, responseBody["can_edit"])
		})

		t.Run("family member can view but not edit", func(t *testing.T) {
			status, responseBody := doJSON(t, client, http.MethodGet, url, familyMemberId, nil)
			require.Equal(t, http.StatusOK, status)
			require.Equal(t, false, responseBody["can_edit"])
		})

		t.Run("outsider cannot view", func(t *testing.T) {
			status, _ := doJSON(t, client, http.MethodGet, url, outsiderId, nil)
			require.Equal(t, http.StatusUnauthorized, status)
		})
	}
}

func UpdateRecipeForbiddenTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d", dbAdapterUri, recipeId)

		status, _ := doJSON(t, client, http.MethodPatch, url, familyMemberId, update_recipe)
		require.Equal(t, http.StatusUnauthorized, status)
	}
}

func MissingRecipeCommentTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl // test boilerplate
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/comment", dbAdapterUri, recipeId)

		t.Run("update is 404", func(t *testing.T) {
			status, _ := doJSON(t, client, http.MethodPatch, url, familyMemberId, editRecipeComment)
			require.Equal(t, http.StatusNotFound, status)
		})

		t.Run("delete is 404", func(t *testing.T) {
			status, _ := doJSON(t, client, http.MethodDelete, url, familyMemberId, nil)
			require.Equal(t, http.StatusNotFound, status)
		})
	}
}

func SoftDeleteRecipeForbiddenTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d", dbAdapterUri, recipeId)

		status, _ := doJSON(t, client, http.MethodDelete, url, familyMemberId, nil)
		require.Equal(t, http.StatusUnauthorized, status)
	}
}

func HardDeleteLiveRecipeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d/hard-delete", dbAdapterUri, recipeId)

		status, _ := doJSON(t, client, http.MethodDelete, url, 1, nil)
		require.Equal(t, http.StatusConflict, status)

		status, _ = doJSON(t, client, http.MethodGet, fmt.Sprintf("%s/recipe/%d", dbAdapterUri, recipeId), 1, nil)
		require.Equal(t, http.StatusOK, status, "the live recipe must be untouched")
	}
}

func SoftDeleteRecipeTwiceTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := fmt.Sprintf("%s/recipe/%d", dbAdapterUri, recipeId)

		status, _ := doJSON(t, client, http.MethodDelete, url, 1, nil)
		require.Equal(t, http.StatusNotFound, status)
	}
}
