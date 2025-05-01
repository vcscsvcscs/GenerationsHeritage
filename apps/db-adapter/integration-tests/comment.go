package integration_tests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

//go:embed payloads/comment.json
var comment []byte

func CommentOnPersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/comment/8"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(comment))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "6")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func GetCommentsOnPersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		url := dbAdapterUri + "/comment/8"

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("X-User-ID", "6")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "StartElementId")
	}
}

//go:embed payloads/edit_comment.json
var edit_comment []byte

func PatchCommentOnPersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		patchURL := dbAdapterUri + "/comment/8"

		patchReq, err := http.NewRequestWithContext(t.Context(), http.MethodPatch, patchURL, bytes.NewBuffer(edit_comment))
		require.NoError(t, err)
		patchReq.Header.Set("Content-Type", "application/json")
		patchReq.Header.Set("X-User-ID", "6")

		patchResp, err := client.Do(patchReq)
		require.NoError(t, err)
		defer patchResp.Body.Close()

		require.Equal(t, http.StatusOK, patchResp.StatusCode)

		var responseBody api.Messages
		err = json.NewDecoder(patchResp.Body).Decode(&responseBody)
		require.NoError(t, err)
		require.NotEmpty(t, responseBody.Comments)
		require.Len(t, *responseBody.Comments, 1)
	}
}

func DeleteCommentOnPersonTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		deleteURL := dbAdapterUri + "/comment/8"
		deleteReq, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, deleteURL, http.NoBody)
		require.NoError(t, err)
		deleteReq.Header.Set("X-User-ID", "6")

		deleteResp, err := client.Do(deleteReq)
		require.NoError(t, err)
		defer deleteResp.Body.Close()

		require.Equal(t, http.StatusOK, deleteResp.StatusCode)
	}
}
