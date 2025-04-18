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

func (srv *server) CreateRelationship(c *gin.Context, params api.CreateRelationshipParams) {
	var relationship *api.CreateRelationshipJSONRequestBody
	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()

	if err := auth.CouldManagePersonUnknownAdmin(actx, session, *relationship.Id1, params.XUserID); err != nil {
		if err := auth.CouldManagePersonUnknownAdmin(actx, session, *relationship.Id1, params.XUserID); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

			return
		}
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	var relationShipResultRaw any
	var relationshipError error
	switch *relationship.Type {
	case api.CreateRelationshipJSONBodyTypeChild:
		relationShipResultRaw, relationshipError = session.ExecuteWrite(
			qctx, memgraph.CreateChildParentRelationship(
				qctx, *relationship.Id1, *relationship.Id2, *relationship.Relationship,
			),
		)
	case api.CreateRelationshipJSONBodyTypeParent:
		relationShipResultRaw, relationshipError = session.ExecuteWrite(
			qctx, memgraph.CreateChildParentRelationship(
				qctx, *relationship.Id1, *relationship.Id2, *relationship.Relationship,
			),
		)
	case api.CreateRelationshipJSONBodyTypeSibling:
		relationShipResultRaw, relationshipError = session.ExecuteWrite(
			qctx, memgraph.CreateSiblingRelationship(
				qctx, *relationship.Id1, *relationship.Id2, *relationship.Relationship,
			),
		)
	case api.CreateRelationshipJSONBodyTypeSpouse:
		relationShipResultRaw, relationshipError = session.ExecuteWrite(
			qctx, memgraph.CreateSpouseRelationship(
				qctx, *relationship.Id1, *relationship.Id2, *relationship.Relationship,
			),
		)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid relationship type"})
	}
	if relationshipError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": relationshipError.Error()})

		return
	}

	c.JSON(http.StatusOK, relationShipResultRaw)
}

func (srv *server) UpdateRelationship(c *gin.Context, id1, id2 int, params api.UpdateRelationshipParams) {
	var relationship *api.UpdateRelationshipJSONBody
	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()

	if err := auth.CouldManagePersonUnknownAdmin(actx, session, *relationship.Id1, params.XUserID); err != nil {
		if err := auth.CouldManagePersonUnknownAdmin(actx, session, *relationship.Id1, params.XUserID); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

			return
		}
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()
	res, err := session.ExecuteWrite(qctx, memgraph.UpdateRelationship(qctx, id1, id2, *relationship.Relationship))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) GetRelationship(c *gin.Context, id1, id2 int, params api.GetRelationshipParams) {
	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()

	if err := auth.CouldManagePersonUnknownAdmin(actx, session, id1, params.XUserID); err != nil {
		if err := auth.CouldManagePersonUnknownAdmin(actx, session, id2, params.XUserID); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

			return
		}
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	res, err := session.ExecuteRead(qctx, memgraph.GetRelationship(qctx, id1, id2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *server) DeleteRelationship(c *gin.Context, id1, id2 int, params api.DeleteRelationshipParams) {
	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	actx, aCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer aCancel()

	if err := auth.CouldManagePersonUnknownAdmin(actx, session, id1, params.XUserID); err != nil {
		if err := auth.CouldManagePersonUnknownAdmin(actx, session, id2, params.XUserID); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": fmt.Sprint("User does not have access to this person", err.Error())})

			return
		}
	}

	qctx, qCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer qCancel()

	_, err := session.ExecuteWrite(qctx, memgraph.DeleteRelationship(qctx, id1, id2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "relationship deleted"})
}
