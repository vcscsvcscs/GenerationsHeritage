package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) DeleteRecipeRelationship(c *gin.Context, recipeId int, params api.DeleteRecipeRelationshipParams) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}

func (srv *server) CreateRecipeRelationship(c *gin.Context, recipeId int, params api.CreateRecipeRelationshipParams) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}
