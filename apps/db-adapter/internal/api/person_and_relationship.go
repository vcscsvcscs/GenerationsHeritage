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

	ctx, cancel := context.WithTimeout(c.Request.Context(), srv.dbOpTimeout)
	defer cancel()

	session := srv.db.NewSession(ctx, neo4j.SessionConfig{})
	defer closeSession(c.Request.Context(), session, srv.dbOpTimeout)

	qctx, qCancel := context.WithTimeout(ctx, srv.dbOpTimeout)
	defer qCancel()

	res, err := session.ExecuteWrite(qctx, memgraph.CreatePerson(qctx, &api.PersonProperties{
		FirstName:        &requestBody.Person.FirstName,
		LastName:         &requestBody.Person.LastName,
		Born:             &requestBody.Person.Born,
		MothersFirstName: &requestBody.Person.MothersFirstName,
		MothersLastName:  &requestBody.Person.MothersLastName,
		Limit:            &requestBody.Person.Limit,
	}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	resMap, ok := res.(map[string]any)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "unexpected result type"})

		return
	}

	relCtx, relCCancel := context.WithTimeout(ctx, srv.dbOpTimeout)
	defer relCCancel()

	personID := resMap["id"].(int)
	var relationShipResultRaw any
	var relationshipError error
	switch *requestBody.Type {
	case api.CreatePersonAndRelationshipJSONBodyTypeChild:
		relationShipResultRaw, relationshipError = session.ExecuteRead(relCtx, memgraph.CreateChildParentRelationship(qctx, personID, id))
	case api.CreatePersonAndRelationshipJSONBodyTypeParent:
		relationShipResultRaw, relationshipError = session.ExecuteRead(relCtx, memgraph.CreateChildParentRelationship(qctx, id, personID))
	case api.CreatePersonAndRelationshipJSONBodyTypeSibling:
		relationShipResultRaw, relationshipError = session.ExecuteRead(relCtx, memgraph.CreateSiblingRelationship(qctx, personID, id))
	case api.CreatePersonAndRelationshipJSONBodyTypeSpouse:
		relationShipResultRaw, relationshipError = session.ExecuteRead(relCtx, memgraph.CreateSpouseRelationship(qctx, personID, id))
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
	}{Person: resMap, Rel: relationShipResultRaw})
}
