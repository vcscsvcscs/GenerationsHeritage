package memgraph

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

func (rp *RelationshipAndPerson) CreateRelationshipAndPerson(driver neo4j.DriverWithContext) (*neo4j.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	if err := rp.Verify(); err != nil {
		return nil, err
	}

	rp.Person.ID = strings.ReplaceAll(uuid.New().String(), "-", "")

	query := fmt.Sprintf(`MATCH (a:Person) WHERE a.id = '%s'`, rp.Relationship.FirstPersonID)

	query = fmt.Sprintf("%s CREATE (b:Person {%s})", query, rp.Person.ToString())

	if rp.Relationship.Direction == "->" {
		query = fmt.Sprintf(`%s CREATE (a)-[r:%s {verified: True}]->(b) RETURN r;`, query, rp.Relationship.Relationship)
	} else if rp.Relationship.Direction == "<-" {
		query = fmt.Sprintf(`%s CREATE (a)<-[r:%s {verified: True}]-(b) RETURN r;`, query, rp.Relationship.Relationship)
	} else {
		query = fmt.Sprintf(`%s CREATE (a)<-[r1:%s {verified: True}]-(b) CREATE (a)-[r2:%s {verified: True}]->(b) RETURN r1, r2, b;`,
			query, rp.Relationship.Relationship, rp.Relationship.Relationship)
	}

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	return result.Single(ctx)
}
