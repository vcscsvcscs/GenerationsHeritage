package memgraph

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

func (r *Relationship) CreateRelationship(driver neo4j.DriverWithContext) (*neo4j.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := r.Verify(); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`MATCH (a:Person), (b:Person)
	WHERE a.ID = %s AND b.ID = %s`, r.FirstPersonID, r.SecondPersonID)

	if r.Direction == "->" {
		query = fmt.Sprintf(`%s CREATE (a)-[r:%s {verified: false}]->(b) RETURN r;`, query, r.Relationship)
	} else if r.Direction == "<-" {
		query = fmt.Sprintf(`%s CREATE (a)<-[r:%s {verified: false}]-(b) RETURN r;`, query, r.Relationship)
	} else {
		query = fmt.Sprintf(`%s CREATE (a)-[r:%s {verified: false}]-(b) RETURN r;`, query, r.Relationship)
	}

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	return result.Single(ctx)
}
