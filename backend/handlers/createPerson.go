package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/backend/memgraph"
)

func CreatePerson(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close(ctx)

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

		person.ID = uuid.New().String()
		if person.Verify() != nil {
			log.Printf("ip: %s error: %s", c.ClientIP(), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person information"})

			return
		}

		query := fmt.Sprintf("CREATE (n:Person {%s}) RETURN n;", person.ToString())

		result, err := session.Run(ctx, query, nil)
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

		c.JSON(200, gin.H{"person": rec.AsMap()})
	}
}
