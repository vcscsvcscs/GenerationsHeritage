package api

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func closeSession(ctx context.Context, session neo4j.SessionWithContext, timeOut time.Duration) {
	sctx, cancel := context.WithTimeout(ctx, timeOut)
	session.Close(sctx)
	cancel()
}
