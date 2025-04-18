package memgraph

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

func InitDatabase(logger *zap.Logger, dbURI, dbUser, dbPassword string) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(dbURI, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		logger.Panic("Unable to init driver due to:", zap.Error(err))
	}

	ctx := context.Background()

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		logger.Panic("Unable to connect to db:", zap.Error(err),
			zap.String("dbUser", dbUser),
			zap.String("dbURI", dbURI),
		)
	}

	return driver
}
