package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) GetFamilyTreeById(c *gin.Context, params api.GetFamilyTreeByIdParams) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetFamilyTreeById(qctx, params.XUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) GetFamilyTreeWithSpousesById(c *gin.Context, params api.GetFamilyTreeWithSpousesByIdParams) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetFamilyTreeWithSpousesById(qctx, params.XUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
