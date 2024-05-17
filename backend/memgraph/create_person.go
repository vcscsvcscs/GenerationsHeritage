package memgraph

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

func (p *Person) CreatePerson(driver neo4j.DriverWithContext) (*neo4j.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := p.Verify(); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("CREATE (n:Person {%s}) RETURN n;", p.ToString())

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	return result.Single(ctx)
}
