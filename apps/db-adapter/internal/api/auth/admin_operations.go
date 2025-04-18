package auth

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
)

// CouldManagePerson determines if an admin has the authority to manage a person.
// It checks if the provided adminId matches the XUserID, and if not, delegates
// the check to CouldManagePersonUnknownAdmin.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - session: The Neo4j session used for database operations.
//   - userId: The ID of the user being managed.
//   - adminId: The ID of the admin attempting to manage the user.
//   - XUserID: The ID of the currently authenticated user.
//
// Returns:
//   - An error if the admin does not have the authority to manage the person,
//     or nil if the operation is allowed.
func CouldManagePerson(ctx context.Context, session neo4j.SessionWithContext, userId, adminId, XUserID int) error {
	if adminId == XUserID {
		return nil
	}

	return CouldManagePersonUnknownAdmin(ctx, session, userId, XUserID)
}

// CouldManagePersonUnknownAdmin checks if a user can manage another person
// when the user is not an admin. It verifies if the provided userId matches
// the XUserID, and if not, it attempts to read the admin relationship between
// the two users from the database.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - session: The Neo4j session used to execute the database query.
//   - userId: The ID of the user attempting to manage another person.
//   - XUserID: The ID of the person being managed.
//
// Returns:
//   - An error if the user is not allowed to manage the person or if there is
//     an issue querying the database. Returns nil if the user is allowed.
func CouldManagePersonUnknownAdmin(ctx context.Context, session neo4j.SessionWithContext, userId, XUserID int) error {
	if userId == XUserID {
		return nil
	}

	_, err := session.ExecuteRead(ctx, memgraph.GetAdminRelationship(ctx, userId, XUserID))
	return err
}
