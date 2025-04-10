package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/api/auth"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) CreateAdminRelationship(c *gin.Context, id1 int, id2 int, params api.CreateAdminRelationshipParams) {
	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()
	if err := auth.CouldManagePerson(actx, session, id1, id2, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have permissions to manage this person with error:", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	res, err := session.ExecuteWrite(qctx, memgraph.CreateAdminRelationship(qctx, id1, id2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) DeleteAdminRelationship(c *gin.Context, id1 int, id2 int, params api.DeleteAdminRelationshipParams) {
	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()
	if err := auth.CouldManagePerson(actx, session, id1, id2, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have permissions to manage this person with error:", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	_, err := session.ExecuteWrite(qctx, memgraph.DeleteAdminRelationship(qctx, id1, id2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "admin relationship was deleted"})
}

func (srv *server) GetAdminRelationship(c *gin.Context, id1 int, id2 int, params api.GetAdminRelationshipParams) {
	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()
	if err := auth.CouldSeePersonsProfile(actx, session, id1, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have permissions to see this person", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	res, err := session.ExecuteRead(qctx, memgraph.GetAdminRelationship(qctx, id1, id2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) GetProfileAdmins(c *gin.Context, id int, params api.GetProfileAdminsParams) {
	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()
	if err := auth.CouldSeePersonsProfile(actx, session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have permissions to see this person", err.Error())})

		return
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	res, err := session.ExecuteRead(qctx, memgraph.GetProfileAdmins(qctx, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) GetManagedProfiles(c *gin.Context, params api.GetManagedProfilesParams) {
	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	res, err := session.ExecuteRead(qctx, memgraph.GetManagedProfiles(qctx, params.XUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}
