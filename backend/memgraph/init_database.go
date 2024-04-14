package memgraph

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitDatabase(dbUri, dbUser, dbPassword string) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		log.Panicln(err)
	}

	ctx := context.Background()
	defer driver.Close(ctx)

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
