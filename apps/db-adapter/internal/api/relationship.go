package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) CreateRelationship(c *gin.Context, params api.CreateRelationshipParams) {}

func (srv *server) GetRelationship(c *gin.Context, id1 int, id2 int, params api.GetRelationshipParams) {
}
