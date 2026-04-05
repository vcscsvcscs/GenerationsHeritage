package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/api/auth"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) GetRecipesByPersonId(
	c *gin.Context, id int, params api.GetRecipesByPersonIdParams,
) {
	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	if err := auth.CouldSeePersonsProfile(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetRecipesByPersonId(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) CreateRecipeForPerson(
	c *gin.Context, id int, params api.CreateRecipeForPersonParams,
) {
	var body api.CreateRecipeForPersonJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	if err := auth.CouldSeePersonsProfile(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.CreateRecipeForPerson(
		qctx, id, &body.Recipe, body.Relationship,
	))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}
