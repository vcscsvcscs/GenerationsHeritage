package memgraph

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

func (r *Relationship) VerifyRelationship(driver neo4j.DriverWithContext) (*neo4j.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := r.Verify(); err != nil {
		return nil, err
	}

	query := ""
	if r.Direction == "->" {
		query = fmt.Sprintf(`MATCH (a)-[r:%s]->(b)`, r.Relationship)
	} else if r.Direction == "<-" {
		query = fmt.Sprintf(`MATCH (a)<-[r:%s]-(b)`, r.Relationship)
	} else {
		query = fmt.Sprintf(`MATCH (a)-[r:%s]-(b)`, r.Relationship)
	}

	query = fmt.Sprintf(`%s WHERE a.ID = %s AND b.ID = %s set r.verified = true return r;`, query, r.FirstPersonID, r.SecondPersonID)

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	return result.Single(ctx)
}
