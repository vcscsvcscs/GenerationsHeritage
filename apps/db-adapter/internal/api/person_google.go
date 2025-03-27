package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (srv *server) GetPersonByGoogleId(c *gin.Context, googleId string) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})
	session.ExecuteRead()
}

func (srv *server) CreatePersonByGoogleIdAndInviteCode(c *gin.Context, googleId string) {}

func (srv *server) CreatePersonByGoogleId(c *gin.Context, googleId string) {}
