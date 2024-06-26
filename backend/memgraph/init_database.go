package memgraph

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitDatabase(dbURI, dbUser, dbPassword string) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(dbURI, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		log.Panicln(err)
	}

	ctx := context.Background()

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Panicln(err)
	}

	if err := createIndexes(driver); err != nil {
		log.Panicln(err)
	}

	if err := createConstraints(driver); err != nil {
		log.Panicln(err)
	}

	return driver
}
