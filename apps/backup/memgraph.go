package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type memgraphDB struct {
	driver neo4j.DriverWithContext
}

func newMemgraphDB(ctx context.Context, c *config) (*memgraphDB, error) {
	d, err := neo4j.NewDriverWithContext(c.memgraphURI, neo4j.BasicAuth(c.memgraphUser, c.memgraphPass, ""))
	if err != nil {
		return nil, fmt.Errorf("memgraph driver: %w", err)
	}

	if err = d.VerifyConnectivity(ctx); err != nil {
		_ = d.Close(ctx)

		return nil, fmt.Errorf("memgraph connectivity: %w", err)
	}

	return &memgraphDB{driver: d}, nil
}

func (m *memgraphDB) Close(ctx context.Context) error {
	return m.driver.Close(ctx) //nolint:wrapcheck // thin wrapper
}

// Dump streams the result of DUMP DATABASE, which yields one Cypher statement per record.
func (m *memgraphDB) Dump(ctx context.Context, emit func(stmt string) error) error {
	s := m.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer s.Close(ctx)

	res, err := s.Run(ctx, "DUMP DATABASE", nil)
	if err != nil {
		return fmt.Errorf("DUMP DATABASE: %w", err)
	}

	for res.Next(ctx) {
		stmt, ok := res.Record().Values[0].(string)
		if !ok {
			return fmt.Errorf("unexpected DUMP DATABASE value %T", res.Record().Values[0])
		}

		if eerr := emit(stmt); eerr != nil {
			return eerr
		}
	}

	if err = res.Err(); err != nil {
		return fmt.Errorf("DUMP DATABASE: %w", err)
	}

	return nil
}

// Exec runs one statement in its own auto-commit transaction, as index DDL cannot run in explicit ones.
func (m *memgraphDB) Exec(ctx context.Context, stmt string) error {
	s := m.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer s.Close(ctx)

	res, err := s.Run(ctx, stmt, nil)
	if err != nil {
		return err //nolint:wrapcheck // wrapped by caller
	}

	_, err = res.Consume(ctx)

	return err //nolint:wrapcheck // wrapped by caller
}

func (m *memgraphDB) IsEmpty(ctx context.Context) (bool, error) {
	s := m.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer s.Close(ctx)

	res, err := s.Run(ctx, "MATCH (n) RETURN 1 LIMIT 1", nil)
	if err != nil {
		return false, err //nolint:wrapcheck // wrapped by caller
	}

	if res.Next(ctx) {
		return false, nil
	}

	return true, res.Err() //nolint:wrapcheck // wrapped by caller
}

func (m *memgraphDB) Wipe(ctx context.Context) error {
	return m.Exec(ctx, "MATCH (n) DETACH DELETE n")
}
