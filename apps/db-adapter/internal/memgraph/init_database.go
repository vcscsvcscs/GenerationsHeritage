package memgraph

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

const (
	databaseReconnectTimeout  = 5 * time.Second
	databaseReconnectAttempts = 5
)

func InitDatabase(logger *zap.Logger, dbURI, dbUser, dbPassword string) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(dbURI, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		logger.Panic("Unable to init driver due to:", zap.Error(err))
	}

	ctx := context.Background()
	for attempt := 1; attempt <= databaseReconnectAttempts; attempt++ {
		err = driver.VerifyConnectivity(ctx)
		if err == nil {
			break
		}

		logger.Warn("Retrying connection to db...",
			zap.Error(err),
			zap.String("dbUser", dbUser),
			zap.String("dbURI", dbURI),
			zap.Int("attempt", attempt),
		)
		time.Sleep(time.Duration(attempt) * databaseReconnectTimeout)
	}

	if err != nil {
		logger.Panic("Unable to connect to db after retries:", zap.Error(err),
			zap.String("dbUser", dbUser),
			zap.String("dbURI", dbURI),
		)
	}

	return driver
}
