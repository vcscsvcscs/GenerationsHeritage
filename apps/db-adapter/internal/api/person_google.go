package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) GetPersonByGoogleId(c *gin.Context, googleId string) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetPersonByGoogleId(qctx, googleId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) CreatePersonByGoogleIdAndInviteCode(c *gin.Context, googleId string) {
	var person struct {
		InviteCode string                `json:"invite_code"`
		Props      *api.PersonProperties `json:"person"`
	}
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	emptyString := ""
	person.Props.InviteCode = &emptyString

	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.UpdatePersonByInviteCode(qctx, person.InviteCode, person.Props))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) CreatePersonByGoogleId(c *gin.Context, googleId string) {
	var person *api.PersonProperties
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}
	person.GoogleId = &googleId // just making sure :)

	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.CreatePerson(qctx, person))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}
