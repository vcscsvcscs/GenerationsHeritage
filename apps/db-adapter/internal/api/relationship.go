package api

import "github.com/gin-gonic/gin"

func (srv *server) CreateRelationship(c *gin.Context, params CreateRelationshipParams) {}

func (srv *server) GetRelationship(c *gin.Context, id1 int, id2 int, params GetRelationshipParams) {}
