package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) CreatePersonAndRelationship(c *gin.Context, id int, params api.CreatePersonAndRelationshipParams) {
	var requestBody api.CreatePersonAndRelationshipJSONBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})

		return
	}

	session := srv.db.NewSession(c.Request.Context(), neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	trs, err := session.BeginTransaction(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})

		return
	}
	defer func() {
		trs.Commit(c.Request.Context())
		trs.Close(c.Request.Context())
	}()

	qctx, qCancel := context.WithTimeout(context.Background(), srv.dbOpTimeout)
	defer qCancel()
	res, err := trs.Run(qctx, memgraph.CreatePersonCypherQuery, map[string]any{
		"Person": requestBody.Person,
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
	personID, ok := singleRes.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Person ID not found in response"})

		return
	}

	actx, acancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer acancel()
	_, aErr := trs.Run(actx, memgraph.CreateAdminRelationshipCypherQuery, map[string]any{
		"id2": personID.(int),
		"id1": params.XUserID,
	})
	if aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": aErr.Error()})

		return
	}

	relCtx, relCCancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer relCCancel()

	var relationShipResultRaw any
	var relationshipError error
	switch *requestBody.Type {
	case api.CreatePersonAndRelationshipJSONBodyTypeChild:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateChildParentRelationshipCypherQuery, map[string]any{
			"childId":            personID.(int),
			"parentId":           id,
			"childRelationship":  requestBody.Relationship,
			"parentRelationship": requestBody.Relationship,
		})
	case api.CreatePersonAndRelationshipJSONBodyTypeParent:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateChildParentRelationshipCypherQuery, map[string]any{
			"childId":            id,
			"parentId":           personID.(int),
			"childRelationship":  requestBody.Relationship,
			"parentRelationship": requestBody.Relationship,
		})
	case api.CreatePersonAndRelationshipJSONBodyTypeSibling:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateSiblingRelationshipCypherQuery, map[string]any{
			"id1":           id,
			"id2":           personID.(int),
			"Relationship1": requestBody.Relationship,
			"Relationship2": requestBody.Relationship,
		})
	case api.CreatePersonAndRelationshipJSONBodyTypeSpouse:
		relationShipResultRaw, relationshipError = trs.Run(relCtx, memgraph.CreateSpouseRelationshipCypherQuery, map[string]any{
			"id1":           personID.(int),
			"id2":           id,
			"Relationship1": requestBody.Relationship,
			"Relationship2": requestBody.Relationship,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid relationship type"})
	}
	if relationshipError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": relationshipError.Error()})

		return
	}

	c.JSON(http.StatusOK, struct {
		Person any `json:"person"`
		Rel    any `json:"relationship"`
	}{Person: singleRes, Rel: relationShipResultRaw})
}
