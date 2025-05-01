package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPersonGoogle runs integration tests for the Person Google API endpoints.
// It requires a running instance of the db-adapter and a valid dbAdapterUri.
// The tests include creating a person by Google ID, and getting a person by Google ID.
// It does not include creating a person by Google ID with an invite code.
func TestPersonGoogle(dbAdapterUri string) func(t *testing.T) {
	return func(t *testing.T) {
		client := &http.Client{}

		t.Run("CreatePersonByGoogleId", createPersonByGoogleIdTest(dbAdapterUri, client))
		t.Run("GetPersonByGoogleId", getPersonByGoogleIdTest(dbAdapterUri, client))
	}
}

func getPersonByGoogleIdTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/google/test-google-id"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)

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

		require.Equal(t, "Alice", responseBody["Props"].(map[string]any)["first_name"])
		require.Equal(t, "Wonderland", responseBody["Props"].(map[string]any)["last_name"])
	}
}

//go:embed payloads/create_person_with_invite_code.json
var create_person_with_invite_code []byte

func CreatePersonByGoogleIdAndInviteCodeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl,lll // won't fix this, as it is a test
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/google/test-google-id"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, url, bytes.NewBuffer(create_person_with_invite_code))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		// t.Log("Response Status Code: ", responseBody)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		_, ok := responseBody["Id"]
		require.True(t, ok)

		require.Equal(t, "John", responseBody["Props"].(map[string]any)["first_name"])
		require.Equal(t, "Doe", responseBody["Props"].(map[string]any)["last_name"])
	}
}

//go:embed payloads/create_person.json
var create_person []byte

func createPersonByGoogleIdTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl,lll // won't fix this, as it is a test
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/google/test-google-id"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(create_person))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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

		require.Equal(t, "Alice", responseBody["Props"].(map[string]any)["first_name"])
		require.Equal(t, "Wonderland", responseBody["Props"].(map[string]any)["last_name"])
	}
}
