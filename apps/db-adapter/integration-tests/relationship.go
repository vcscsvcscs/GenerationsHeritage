package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed payloads/create_relationship_child.json
var create_relationship_child []byte

//go:embed payloads/create_relationship_parent.json
var create_relationship_parent []byte

//go:embed payloads/create_relationship_sibling.json
var create_relationship_sibling []byte

//go:embed payloads/create_relationship_spouse.json
var create_relationship_spouse []byte

func CreateRelationshipsTest(dbAdapterURI string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("CreatePerson6", CreatePersonTest(dbAdapterURI, client))
		t.Run("CreatePerson7", CreatePersonTest(dbAdapterURI, client))
		t.Run("CreatePerson8", CreatePersonTest(dbAdapterURI, client))
		t.Run("CreateChildRelationship", CreateRelationshipTest(dbAdapterURI, &create_relationship_child, client))
		t.Run("CreateParentRelationship", CreateRelationshipTest(dbAdapterURI, &create_relationship_parent, client))
		t.Run("CreateSiblingRelationship", CreateRelationshipTest(dbAdapterURI, &create_relationship_sibling, client))
		t.Run("CreateSpouseRelationship", CreateRelationshipTest(dbAdapterURI, &create_relationship_spouse, client))
	}
}

func CreateRelationshipTest(dbAdapterURI string, payload *[]byte, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterURI + "/relationship"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(*payload))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "0")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody []any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

//go:embed payloads/verify_relationship.json
var verify_relationship []byte

func UpdateRelationshipTest(dbAdapterURI string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterURI + "/relationship/6/7"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, url, bytes.NewBuffer(verify_relationship))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "7")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody []any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func GetRelationshipTest(dbAdapterURI string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterURI + "/relationship/5/8"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "5")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		_, ok := responseBody["Id"]
		require.True(t, ok)

		require.Equal(t, "5", responseBody["StartElementId"])
		require.Equal(t, "8", responseBody["EndElementId"])
	}
}

func DeleteRelationshipTest(dbAdapterURI string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterURI + "/relationship/5/8"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "5")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)
		_, ok := responseBody["msg"]
		require.True(t, ok)
	}
}
