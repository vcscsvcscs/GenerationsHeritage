package auth

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CouldSeePersonsProfile(ctx context.Context, session neo4j.SessionWithContext, userId, XUserID int) error {
	if CouldManagePersonUnknownAdmin(ctx, session, userId, XUserID) == nil {
		return nil
	}

}
