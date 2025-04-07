package auth

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
)

// This function checks if the user has permission to manage another user's profile, it returns an error if the user does not have permission.
func CouldManagePerson(ctx context.Context, session neo4j.SessionWithContext, userId, adminId, XUserID int) error {
	if adminId == XUserID {
		return nil
	}

	return CouldManagePersonUnknownAdmin(ctx, session, userId, XUserID)
}

// This function checks if the user has permission to manage another user's profile, it returns an error if the user does not have permission.
func CouldManagePersonUnknownAdmin(ctx context.Context, session neo4j.SessionWithContext, userId, XUserID int) error {
	if userId == XUserID {
		return nil
	}

	_, err := session.ExecuteRead(ctx, memgraph.GetAdminRelationship(ctx, userId, XUserID), nil)
	return err
}
