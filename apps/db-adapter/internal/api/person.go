package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) CreatePerson(c *gin.Context, params api.CreatePersonParams) {
	var person *api.PersonProperties
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	adminList := []struct {
		Id   *int    "json:\"id,omitempty\""
		Name *string "json:\"name,omitempty\""
	}{
		{Id: &params.XUserID, Name: &params.XUserName},
	}
	person.AllowAdminAccess = &adminList

	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.CreatePerson(qctx, person))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) GetPersonById(c *gin.Context, id int, params api.GetPersonByIdParams) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer acancel()
	if userWithIdHasAccessToGivenPerson(actx, session, params.XUserID, id) == accessModeNone {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User does not have access to this person"})

		return
	}

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteRead(qctx, memgraph.GetPersonById(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) SoftDeletePerson(c *gin.Context, id int, params api.SoftDeletePersonParams) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer acancel()
	if userWithIdHasAccessToGivenPerson(actx, session, params.XUserID, id) != accessModeWrite {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User does not have access to this person"})

		return
	}

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.SoftDeletePerson(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) UpdatePerson(c *gin.Context, id int, params api.UpdatePersonParams) {
	var person *api.PersonProperties
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	actx, acancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer acancel()
	session := srv.db.NewSession(actx, neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	if userWithIdHasAccessToGivenPerson(actx, session, params.XUserID, id) != accessModeWrite {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User does not have access to this person"})

		return
	}

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.UpdatePerson(qctx, id, person))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) HardDeletePerson(c *gin.Context, id int, params api.HardDeletePersonParams) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer cancel()
	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer acancel()
	if userWithIdHasAccessToGivenPerson(actx, session, params.XUserID, id) != accessModeWrite {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "User does not have access to this person"})

		return
	}

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.HardDeletePerson(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}
