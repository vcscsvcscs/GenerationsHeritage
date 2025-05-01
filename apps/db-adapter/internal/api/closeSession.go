package api

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

// closeSession closes a Neo4j session with a specified timeout.
// It ensures that the session is properly closed within the given timeout duration.
// If an error occurs during the session closure, it logs the error using the provided logger.
//
// Parameters:
//   - ctx: The parent context for managing the session closure.
//   - logger: The logger instance used to log any errors during session closure.
//   - session: The Neo4j session to be closed.
//   - timeOut: The maximum duration allowed for closing the session.
func closeSession(ctx context.Context, logger *zap.Logger, session neo4j.SessionWithContext, timeOut time.Duration) {
	sctx, cancel := context.WithTimeout(ctx, timeOut)
	if err := session.Close(sctx); err != nil {
		logger.Error("Error closing session", zap.Error(err))
	}
	cancel()
}
