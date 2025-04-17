package api

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

func closeSession(ctx context.Context, logger *zap.Logger, session neo4j.SessionWithContext, timeOut time.Duration) {
	sctx, cancel := context.WithTimeout(ctx, timeOut)
	if err := session.Close(sctx); err != nil {
		logger.Error("Error closing session", zap.Error(err))
	}
	cancel()
}
