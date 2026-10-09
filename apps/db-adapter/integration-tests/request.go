package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// familyMemberId is a person directly related to person 1 who is neither person 1 nor an admin of it.
	familyMemberId = 2
	// outsiderId is a person that is not part of person 1's family.
	outsiderId = 6
	// current unix time is between these bounds in seconds but far above them in milliseconds.
	minUnixSeconds = 1e9
	maxUnixSeconds = 1e11
)

// doJSON sends a request as userId and returns the status code and the decoded JSON object body.
func doJSON(t *testing.T, client *http.Client, method, url string, userId int, body []byte) (status int, responseBody map[string]any) {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), method, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", strconv.Itoa(userId))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&responseBody))

	return resp.StatusCode, responseBody
}

// requireUnixSeconds asserts that value is a unix timestamp in seconds (milliseconds would be 1000 times larger).
func requireUnixSeconds(t *testing.T, value any) {
	t.Helper()

	timestamp, ok := value.(float64)
	require.True(t, ok, "timestamp should be a number, got %v", value)
	require.Greater(t, timestamp, minUnixSeconds)
	require.Less(t, timestamp, maxUnixSeconds)
}

// entryForRecipe returns the entry whose recipe has the given id, or nil.
func entryForRecipe(entries []any, id int) map[string]any {
	for _, raw := range entries {
		entry, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		recipe, ok := entry["recipe"].(map[string]any)
		if ok && int(recipe["Id"].(float64)) == id {
			return entry
		}
	}

	return nil
}
