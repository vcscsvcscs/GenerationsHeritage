package integration_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetFamilyTreeByIdTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl,lll // won't fix this, as it is a test
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

		require.Len(t, responseBody.People, 4)        //nolint:mnd // 4 people in the family tree
		require.Len(t, responseBody.Relationships, 4) //nolint:mnd // 4 relationships in the family tree
	}
}

func GetFamilyTreeWithSpousesByIdTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl,lll // won't fix this, as it is a test
	return func(t *testing.T) {
		url := dbAdapterUri + "/family-tree-with-spouses"

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

		require.Len(t, responseBody.People, 5)        //nolint:mnd // 5 people in the family tree
		require.Len(t, responseBody.Relationships, 5) //nolint:mnd // 5 relationships in the family tree
	}
}
