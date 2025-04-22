package integration_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetFamilyTreeByIdTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/family-tree"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody struct {
			People        []any `json:"people"`
			Relationships []any `json:"relationships"`
		}

		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, 4, len(responseBody.People))
		require.Equal(t, 4, len(responseBody.Relationships))
	}
}

func GetFamilyTreeWithSpousesByIdTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/family-tree-with-spouses"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "3")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody struct {
			People        []any `json:"people"`
			Relationships []any `json:"relationships"`
		}

		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		require.Equal(t, 5, len(responseBody.People))
		require.Equal(t, 6, len(responseBody.Relationships))
	}
}
