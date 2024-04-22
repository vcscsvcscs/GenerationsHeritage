package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/backend/memgraph"
)

func UpdatePerson(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil || c.ContentType() != "application/json" {
			log.Printf("ip: %s error: request body is empty or content type is not application/json", c.ClientIP())
			c.JSON(http.StatusBadRequest, gin.H{"error": "content type must be application/json and request body must not be empty"})

			return
		}

		var person memgraph.Person
		err := json.NewDecoder(c.Request.Body).Decode(&person)
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})

			return
		}

		if person.ID == "" {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No ID provided"})

			return
		}

		rec, err := person.UpdatePerson(driver)
		if err != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusNotFound, gin.H{"error": "could not update person with information provided"})

			return
		}

		c.JSON(http.StatusOK, gin.H{"person": rec.AsMap()})
	}
}
