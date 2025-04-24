package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/api/auth"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
	"go.uber.org/zap"
)

func (srv *server) CreatePerson(c *gin.Context, params api.CreatePersonParams) {
	var person *api.PersonProperties
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	trs, err := session.BeginTransaction(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}
	defer func() {
		if err := trs.Close(c.Request.Context()); err != nil { //nolint:govet // ignore shadowing
			srv.logger.Error("failed to close transaction", zap.Error(err))
		}
	}()

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	convertedPerson := memgraph.StructToMap(person)
	res, err := trs.Run(qctx, memgraph.CreatePersonCypherQuery, map[string]any{
		"Person": convertedPerson,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	singleRes, err := res.Single(qctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	createdPerson, ok := singleRes.Get("person")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Person not found in db response"})

		return
	}

	personId := createdPerson.(dbtype.Node).Id //nolint:staticcheck // this is a difference in neo4j and memgraph

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	_, aErr := trs.Run(actx, memgraph.CreateAdminRelationshipCypherQuery, map[string]any{
		"id2": personId,
		"id1": params.XUserID,
	})
	if aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": aErr.Error()})

		return
	}

	if err := trs.Commit(c.Request.Context()); err != nil {
		srv.logger.Error("failed to commit transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, createdPerson)
}

func (srv *server) GetPersonById(c *gin.Context, id int, params api.GetPersonByIdParams) { //nolint:dupl,lll // This just does not worth abstracting anymore
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
	res, err := session.ExecuteRead(qctx, memgraph.GetPersonById(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) SoftDeletePerson(c *gin.Context, id int, params api.SoftDeletePersonParams) { //nolint:dupl,lll // This just does not worth abstracting anymore
	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	if err := auth.CouldManagePersonUnknownAdmin(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	_, err := session.ExecuteWrite(qctx, memgraph.SoftDeletePerson(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Person soft deleted"})
}

func (srv *server) UpdatePerson(c *gin.Context, id int, params api.UpdatePersonParams) {
	var person *api.PersonProperties
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	if err := auth.CouldManagePersonUnknownAdmin(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.UpdatePerson(qctx, id, person))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) HardDeletePerson(c *gin.Context, id int, params api.HardDeletePersonParams) { //nolint:dupl,lll // This just does not worth abstracting anymore
	session := srv.createSessionWithTimeout(c.Request.Context())
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	if err := auth.CouldManagePersonUnknownAdmin(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	_, err := session.ExecuteWrite(qctx, memgraph.HardDeletePerson(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"description": "Person hard deleted",
	})
}
