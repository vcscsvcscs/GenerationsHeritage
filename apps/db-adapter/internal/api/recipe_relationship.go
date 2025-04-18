package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) DeleteRecipeRelationship( //nolint:revive // not implemented
	c *gin.Context, recipeId int, params api.DeleteRecipeRelationshipParams, //nolint:revive // not implemented
) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}

func (srv *server) CreateRecipeRelationship( //nolint:revive // not implemented
	c *gin.Context, recipeId int, params api.CreateRecipeRelationshipParams, //nolint:revive // not implemented
) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}
