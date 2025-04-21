package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

//go:embed payloads/create_other_person.json
var create_other_person []byte

func CreatePersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/person"

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

		require.Equal(t, "Jhon", responseBody["Props"].(map[string]any)["first_name"])
		require.Equal(t, "Doe", responseBody["Props"].(map[string]any)["last_name"])
	}
}

func GetPersonById(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/0"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "0")

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

func SoftDeletePersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl,lll // This just does not worth abstracting anymore
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/0"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "0")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.Equal(t, "Person soft deleted", responseBody["description"])
	}
}

//go:embed payloads/update_person.json
var update_person []byte

func UpdatePersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/1"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, url, bytes.NewBuffer(update_person))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")

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

		require.Equal(t, "John", responseBody["Props"].(map[string]any)["first_name"])
		require.Equal(t, "Doe", responseBody["Props"].(map[string]any)["last_name"])
		require.Equal(t, "ABCD1234", responseBody["Props"].(map[string]any)["invite_code"])
	}
}

//go:embed payloads/generate_invite_code.json
var generate_invite_code []byte

func UpdatePersonWithInviteCodeTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/1"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, url, bytes.NewBuffer(generate_invite_code))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "0")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody api.Person
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotNil(t, responseBody.Id)

		require.Equal(t, "test-invite-code", *responseBody.Props.InviteCode)
	}
}

func HardDeletePersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) { //nolint:dupl,lll // This just does not worth abstracting anymore
	return func(t *testing.T) {
		url := dbAdapterUri + "/person/0/hard-delete"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "0")

		// Send the request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)

		// Validate the response
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.Equal(t, "Person hard deleted", responseBody["description"])
	}
}
