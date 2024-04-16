package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/backend/memgraph"
)

func DeleteRelationship(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var relationship memgraph.Relationship
		if err := c.ShouldBindJSON(&relationship); err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		err := relationship.DeleteRelationship(driver)
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

			return
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "relationship deleted successfully"})
	}
}
