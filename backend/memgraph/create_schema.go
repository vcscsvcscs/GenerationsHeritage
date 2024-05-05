package memgraph

import (
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
		`CREATE INDEX ON :Person(id);`,
		`CREATE INDEX ON :Person(last_name);`,
		`CREATE INDEX ON :Person(first_name);`,
		`CREATE INDEX ON :Person(born);`,
		`CREATE INDEX ON :Person(mothers_firstname);`,
		`CREATE INDEX ON :Person(mothers_last_name);`,
	}

	// Run index queries via implicit auto-commit transaction
	for _, index := range indexes {
		_, err := session.Run(ctx, index, nil)
		if err != nil {
			return err
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
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.id);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.last_name);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.first_name);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.born);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.mothers_firstname);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT EXISTS (n.mothers_lastname);`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT n.id IS UNIQUE;`,
		`CREATE CONSTRAINT ON (n:Person) ASSERT n.last_name, n.first_name, n.born, n.mothers_first_name, n.mothers_last_name IS UNIQUE;`,
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
