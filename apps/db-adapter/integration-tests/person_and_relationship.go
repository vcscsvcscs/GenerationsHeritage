package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed payloads/create_person_and_relationship_child.json
var create_child []byte

//go:embed payloads/create_person_and_relationship_parent.json
var create_parent []byte

//go:embed payloads/create_person_and_relationship_spouse.json
var create_spouse []byte

//go:embed payloads/create_person_and_relationship_sibling.json
var create_sibling []byte

func CreateAFamilyTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("CreateChild", CreatePersonAndRelationshipTest(dbAdapterUri, &create_child, client))
		t.Run("CreateParent", CreatePersonAndRelationshipTest(dbAdapterUri, &create_parent, client))
		t.Run("CreateSpouse", CreatePersonAndRelationshipTest(dbAdapterUri, &create_spouse, client))
		t.Run("CreateSibling", CreatePersonAndRelationshipTest(dbAdapterUri, &create_sibling, client))
	}
}

func CreatePersonAndRelationshipTest(dbAdapterUri string, payload *[]byte, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/person_and_relationship/1"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(create_other_person))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "0")
		req.Header.Set("X-User-Name", "application/json")

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
		_, ok = responseBody["person"]
		require.True(t, ok)
		_, ok = responseBody["relationship"]
		require.True(t, ok)
	}
}
