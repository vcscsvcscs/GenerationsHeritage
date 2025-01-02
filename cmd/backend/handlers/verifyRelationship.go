package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/pkg/memgraph"
)

func VerifyRelationship(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var relationship memgraph.Relationship
		if err := c.ShouldBindJSON(&relationship); err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		rec, err := relationship.VerifyRelationship(driver)
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

			return
		}

		c.JSON(http.StatusOK, gin.H{"relationship": rec.AsMap()})
	}
}
