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

func (srv *server) CommentOnRecipe(c *gin.Context, id int, params api.CommentOnRecipeParams) {
	var body api.CommentOnRecipeJSONRequestBody
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

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.CommentOnRecipe(qctx, params.XUserID, id, body.Message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) GetRecipeComments(c *gin.Context, id int, params api.GetRecipeCommentsParams) {
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
	res, err := session.ExecuteRead(qctx, memgraph.GetRecipeComments(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) UpdateRecipeComment(c *gin.Context, id int, params api.UpdateRecipeCommentParams) {
	var body api.UpdateRecipeCommentJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.UpdateRecipeComment(qctx, params.XUserID, id, body.Message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) DeleteRecipeComment(c *gin.Context, id int, params api.DeleteRecipeCommentParams) {
	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	_, err := session.ExecuteWrite(qctx, memgraph.DeleteRecipeComment(qctx, params.XUserID, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Recipe comment deleted"})
}
