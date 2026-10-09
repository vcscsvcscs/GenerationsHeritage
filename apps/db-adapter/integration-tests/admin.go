package integration_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateAdminRelationship() {
}

func DeleteAdminRelationship() {
}

func GetAdminRelationship() {
}

// GetAdminRelationshipStatusTest checks that a missing admin relationship is 403 and a missing person 404, never 500.
func GetAdminRelationshipStatusTest(dbAdapterUri string, client *http.Client) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("not an admin", func(t *testing.T) {
			url := fmt.Sprintf("%s/admin/1/%d", dbAdapterUri, familyMemberId)

			status, _ := doJSON(t, client, http.MethodGet, url, 1, nil)
			require.Equal(t, http.StatusForbidden, status)
		})

		t.Run("missing person", func(t *testing.T) {
			status, _ := doJSON(t, client, http.MethodGet, dbAdapterUri+"/admin/1/99999", 1, nil)
			require.Equal(t, http.StatusNotFound, status)
		})
	}
}

func GetProfileAdmins() {
}

func GetManagedProfiles() {
}
