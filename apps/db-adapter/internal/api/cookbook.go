package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) GetFamilyCookbook(c *gin.Context, params api.GetFamilyCookbookParams) {
	if params.Distance < 1 || params.Distance > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "distance must be between 1 and 10"})

		return
	}

	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetFamilyCookbook(qctx, params.XUserID, params.Distance))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}
