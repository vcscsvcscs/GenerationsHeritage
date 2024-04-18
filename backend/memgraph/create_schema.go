package memgraph

import (
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

const dbCreateSchemaTimeout = 10 * time.Second

func createIndexes(driver neo4j.DriverWithContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbCreateSchemaTimeout)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	indexes := []string{
		`CREATE INDEX ON :Person(ID);`,
		`CREATE INDEX ON :Person(Lastname);`,
		`CREATE INDEX ON :Person(Firstname);`,
		`CREATE INDEX ON :Person(Born);`,
		`CREATE INDEX ON :Person(MothersFirstname);`,
		`CREATE INDEX ON :Person(MothersLastname);`,
	}

	// Run index queries via implicit auto-commit transaction
	for _, index := range indexes {
		_, err := session.Run(ctx, index, nil)
		if err != nil {
			log.Panicln(err)
		}
	}

	return nil
}

func createConstraints(driver neo4j.DriverWithContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbCreateSchemaTimeout)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	constraints := []string{
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.ID);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.Lastname);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.Firstname);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.Born);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.MothersFirstname);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.MothersLastname);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT n.ID IS UNIQUE;`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT n.Lastname, n.Firstname, n.Born, n.MothersFirstname, n.MothersLastname IS UNIQUE;`,
	}

	// Run index queries via implicit auto-commit transaction
	for _, constraint := range constraints {
		_, err := session.Run(ctx, constraint, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
