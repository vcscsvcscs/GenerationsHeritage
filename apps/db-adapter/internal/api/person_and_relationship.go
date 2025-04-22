package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/api/auth"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
	"go.uber.org/zap"
)

func (srv *server) CreatePersonAndRelationship(c *gin.Context, id int, params api.CreatePersonAndRelationshipParams) {
	var requestBody api.CreatePersonAndRelationshipJSONBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), srv.logger, session, srv.dbOpTimeout)

	if err := auth.CouldManagePersonUnknownAdmin(c.Request.Context(), session, id, params.XUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})

		return
	}

	trs, err := session.BeginTransaction(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}
	defer func() {
		if err := trs.Close(c.Request.Context()); err != nil { //nolint:govet // ignore errcheck
			srv.logger.Error("failed to close transaction", zap.Error(err))
		}
	}()

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	convertedPerson := memgraph.StructToMap(requestBody.Person)
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
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Person ID not found in response"})

		return
	}

	personID := createdPerson.(dbtype.Node).Id //nolint:staticcheck // this is a difference in neo4j and memgraph

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	_, aErr := trs.Run(actx, memgraph.CreateAdminRelationshipCypherQuery, map[string]any{
		"id2": personID,
		"id1": params.XUserID,
	})
	if aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": aErr.Error()})

		return
	}

	relCtx, relCCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer relCCancel()

	convertedRelationship := memgraph.StructToMap(requestBody.Relationship)

	var relationShipResultRaw neo4j.ResultWithContext
	var relationshipError error
	switch *requestBody.Type {
	case api.CreatePersonAndRelationshipJSONBodyTypeChild:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateChildParentRelationshipCypherQuery, map[string]any{
			"childId":            personID,
			"parentId":           id,
			"childRelationship":  convertedRelationship,
			"parentRelationship": convertedRelationship,
		})
	case api.CreatePersonAndRelationshipJSONBodyTypeParent:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateChildParentRelationshipCypherQuery, map[string]any{
			"childId":            id,
			"parentId":           personID,
			"childRelationship":  convertedRelationship,
			"parentRelationship": convertedRelationship,
		})
	case api.CreatePersonAndRelationshipJSONBodyTypeSibling:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateSiblingRelationshipCypherQuery, map[string]any{
			"id1":           id,
			"id2":           personID,
			"Relationship1": convertedRelationship,
			"Relationship2": convertedRelationship,
		})
	case api.CreatePersonAndRelationshipJSONBodyTypeSpouse:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateSpouseRelationshipCypherQuery, map[string]any{
			"id1":           personID,
			"id2":           id,
			"Relationship1": convertedRelationship,
			"Relationship2": convertedRelationship,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid relationship type"})
	}
	if relationshipError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": relationshipError.Error()})

		return
	}

	if err := trs.Commit(c.Request.Context()); err != nil {
		srv.logger.Error("failed to commit transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}

	relationshipsSingle, err := relationShipResultRaw.Single(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "no relationship was created" + err.Error()})

		return
	}

	relationships, ok := relationshipsSingle.Get("relationships")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "no relationship was created"})

		return
	}

	c.JSON(http.StatusOK, struct {
		Person any `json:"person"`
		Rel    any `json:"relationship"`
	}{Person: singleRes, Rel: relationships})
}
