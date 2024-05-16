package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func ViewFamiliyTree(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close(ctx)

		id := c.GetString("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})

			return
		}

		query := `
		MATCH (n:Person {ID: $person_id})-[p:Parent*1..]->(family:Person)
		OPTIONAL MATCH (family)-[c:Child]->(children:Person)
		WITH family, p, children, c, n
		OPTIONAL MATCH (children)<-[p2:Parent]-(OtherParents:Person)
		WITH family, p, children, c, OtherParents, p2,n
		OPTIONAL MATCH (family)-[s:Spouse]-(spouse:Person)
		RETURN family, p, children, c, OtherParents, p2, spouse, s, n;`

		result, err := session.Run(ctx, query, map[string]any{"person_id": id})
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

			return
		}

		rec, err := result.Collect(ctx)
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusNotFound, gin.H{"error": "could not find family tree for person with id: " + id})

			return
		}

		c.JSON(200, rec)
	}
}
