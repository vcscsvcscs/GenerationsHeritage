package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) GetRecipesByPersonId(c *gin.Context, id int, params api.GetRecipesByPersonIdParams) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}
