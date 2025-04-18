package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) GetRecipesByPersonId( //nolint:revive // not implemented
	c *gin.Context, id int, params api.GetRecipesByPersonIdParams, //nolint:revive // not implemented
) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "not implemented"})
}
