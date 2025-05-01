package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) GetFamilyTreeById(c *gin.Context, params api.GetFamilyTreeByIdParams) {
	session := srv.createSessionWithTimeout(c.Request.Context())

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetFamilyTreeById(qctx, params.XUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	results := FamilyTree{
		People:        []any{},
		Relationships: []any{},
	}
	err = FlattenFamilyTree(res, &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (srv *server) GetFamilyTreeWithSpousesById(
	c *gin.Context, params api.GetFamilyTreeWithSpousesByIdParams,
) {
	session := srv.createSessionWithTimeout(c.Request.Context())

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetFamilyTreeWithSpousesById(qctx, params.XUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	results := FamilyTree{
		People:        []any{},
		Relationships: []any{},
	}
	err = FlattenFamilyTree(res, &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, results)
}
