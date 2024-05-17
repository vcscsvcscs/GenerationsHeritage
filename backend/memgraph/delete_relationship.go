package memgraph

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

func (r *Relationship) DeleteRelationship(driver neo4j.DriverWithContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := r.Verify(); err != nil {
		return err
	}

	query := ""
	if r.Direction == "->" {
		query = fmt.Sprintf(`MATCH (a)-[r:%s]->(b)`, r.Relationship)
	} else if r.Direction == "<-" {
		query = fmt.Sprintf(`MATCH (a)<-[r:%s]-(b)`, r.Relationship)
	} else {
		query = fmt.Sprintf(`MATCH (a)-[r:%s]-(b)`, r.Relationship)
	}

	query = fmt.Sprintf(`%s WHERE a.id = '%s' AND b.id = '%s' DELETE r;`, query, r.FirstPersonID, r.SecondPersonID)

	_, err := session.Run(ctx, query, nil)

	return err
}
