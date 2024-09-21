package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func auth(c *gin.Context) {
	id := c.Request.Header.Get("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no id provided"})
		c.Abort()
	}

	c.Set("id", id)
	c.Next()

}
