package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) SoftDeleteRecipe(c *gin.Context, id int, params api.SoftDeleteRecipeParams) { //nolint:revive // not implemented
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}

func (srv *server) UpdateRecipe(c *gin.Context, id int, params api.UpdateRecipeParams) { //nolint:revive // not implemented
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}

func (srv *server) HardDeleteRecipe(c *gin.Context, id int, params api.HardDeleteRecipeParams) { //nolint:revive // not implemented
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}
