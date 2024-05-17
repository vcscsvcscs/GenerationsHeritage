package memgraph

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

func (p *Person) DeletePerson(driver neo4j.DriverWithContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := p.Verify(); err != nil {
		return err
	}

	query := fmt.Sprintf("MATCH (n:Person {id: '%s'}) DELETE n;", p.ID)

	_, err := session.Run(ctx, query, nil)

	return err
}
