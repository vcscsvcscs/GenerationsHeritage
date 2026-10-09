package memgraph

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdminAndPersonLookupsReportNotFound(t *testing.T) {
	t.Run("get admin relationship", func(t *testing.T) {
		params := map[string]any{"id1": 2, "id2": 1}
		_, err := GetAdminRelationship(context.Background(), 1, 2)(txRunning(GetAdminRelationshipCypherQuery, params, emptyResult()))

		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("get admin relationship found", func(t *testing.T) {
		params := map[string]any{"id1": 2, "id2": 1}
		result, err := GetAdminRelationship(context.Background(), 1, 2)(
			txRunning(GetAdminRelationshipCypherQuery, params, recordResult([]string{"relationship"}, "rel")),
		)

		require.NoError(t, err)
		require.Equal(t, map[string]any{"relationship": "rel"}, result)
	})

	t.Run("create admin relationship", func(t *testing.T) {
		params := map[string]any{"id1": 2, "id2": 1}
		_, err := CreateAdminRelationship(context.Background(), 1, 2)(txRunning(CreateAdminRelationshipCypherQuery, params, emptyResult()))

		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("get person", func(t *testing.T) {
		_, err := GetPersonById(context.Background(), 1)(txRunning(GetPersonCypherQuery, map[string]any{"id": 1}, emptyResult()))

		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("get person found", func(t *testing.T) {
		result, err := GetPersonById(context.Background(), 1)(
			txRunning(GetPersonCypherQuery, map[string]any{"id": 1}, recordResult([]string{"person"}, "p")),
		)

		require.NoError(t, err)
		require.Equal(t, "p", result)
	})
}
