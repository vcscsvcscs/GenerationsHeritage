package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestGetAdminRelationshipStatus(t *testing.T) {
	params := api.GetAdminRelationshipParams{XUserID: 1}

	t.Run("existing relationship", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("no relationship between existing persons is 403", func(t *testing.T) {
		srv := recipeServer([]error{memgraph.ErrNotFound, nil, nil}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("missing person is 404", func(t *testing.T) {
		srv := recipeServer([]error{memgraph.ErrNotFound, nil, memgraph.ErrNotFound}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("other failures stay 500", func(t *testing.T) {
		srv := recipeServer([]error{errors.New("db error")}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateAdminRelationshipMissingPerson(t *testing.T) {
	srv := recipeServer(nil, nil, memgraph.ErrNotFound)
	c, w := jsonRequest(t, http.MethodPost, "")

	srv.CreateAdminRelationship(c, 1, 2, api.CreateAdminRelationshipParams{XUserID: 1})

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPersonByIdMissing(t *testing.T) {
	srv := recipeServer([]error{memgraph.ErrNotFound}, nil, nil)
	c, w := jsonRequest(t, http.MethodGet, "")

	srv.GetPersonById(c, 1, api.GetPersonByIdParams{XUserID: 1})

	assert.Equal(t, http.StatusNotFound, w.Code)
}
