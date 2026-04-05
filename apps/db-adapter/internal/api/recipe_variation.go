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

func (srv *server) CreateRecipeVariation(c *gin.Context, id int, params api.CreateRecipeVariationParams) {
	var body api.CreateRecipeVariationJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	if err := auth.CouldSeeRecipe(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this recipe", err.Error())})

		return
	}

	var variationNotes string
	if body.VariationNotes != nil {
		variationNotes = *body.VariationNotes
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.CreateRecipeVariation(
		qctx, id, params.XUserID, &body.Recipe, variationNotes, body.Relationship,
	))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) GetRecipeVariations(c *gin.Context, id int, params api.GetRecipeVariationsParams) {
	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	if err := auth.CouldSeeRecipe(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this recipe", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetRecipeVariations(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}
