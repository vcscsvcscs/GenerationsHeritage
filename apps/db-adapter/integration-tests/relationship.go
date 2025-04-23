package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/stretchr/testify/require"
)

//go:embed payloads/create_relationship_child.json
var create_relationship_child []byte

func CreateRelationshipsTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("CreateChildRelationship", CreateRelationshipTest(dbAdapterUri, &create_relationship_child, client))
	}
}

func CreateRelationshipTest(dbAdapterUri string, payload *[]byte, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/relationship"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(*payload))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")
		req.Header.Set("X-User-Name", "application/json")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody []dbtype.Relationship
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func UpdateRelationship() {
}

func GetRelationship(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/relationship/2/1"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "3")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		_, ok := responseBody["id"]
		require.True(t, ok)

		require.Equal(t, 2, responseBody["Props"].(map[string]any)["start"])
		require.Equal(t, 1, responseBody["Props"].(map[string]any)["end"])
	}
}

func DeleteRelationship() {
}
