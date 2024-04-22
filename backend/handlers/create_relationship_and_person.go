package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/backend/memgraph"
)

func CreateRelationshipAndPerson(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rp memgraph.RelationshipAndPerson
		if err := c.ShouldBindJSON(&rp); err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		rec, err := rp.CreateRelationshipAndPerson(driver)
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

			return
		}

		c.JSON(http.StatusCreated, gin.H{"relationship": rec.AsMap()})
	}
}
