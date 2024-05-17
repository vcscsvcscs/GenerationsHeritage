package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func ViewPerson(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close(ctx)

		id := c.Query("id")
		if id == "" {
			id = c.GetString("id")
		}

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})

			return
		}

		result, err := session.Run(ctx, "MATCH (n:Person) WHERE n.id = $person_id RETURN n;", map[string]any{"person_id": id})
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

			return
		}

		rec, err := result.Single(ctx)
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusNotFound, gin.H{"error": "could not find person with information provided"})

			return
		}

		c.JSON(200, rec.AsMap()["n"])
	}
}
