package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	memgraphMock "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

var errDenied = errors.New("denied")

// recipeServer returns a server whose authorization and read queries consume reads in order
// (nil allows, an error denies); any further read is a test failure. Writes return writeResult, writeErr.
func recipeServer(reads []error, writeResult any, writeErr error) *server {
	mockSession := new(memgraphMock.SessionWithContext)
	mockDriver := new(memgraphMock.DriverWithContext)
	mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)

	for _, readErr := range reads {
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(map[string]any{"read": true}, readErr).Once()
	}

	mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(writeResult, writeErr).Maybe()
	mockSession.On("Close", mock.Anything).Return(nil)

	return &server{db: mockDriver, dbOpTimeout: 2 * time.Second}
}

func jsonRequest(t *testing.T, method, body string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequestWithContext(t.Context(), method, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	return c, w
}

func TestDbErrorStatus(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, dbErrorStatus(memgraph.ErrNotFound))
	assert.Equal(t, http.StatusNotFound, dbErrorStatus(fmt.Errorf("wrapped: %w", memgraph.ErrNotFound)))
	assert.Equal(t, http.StatusConflict, dbErrorStatus(memgraph.ErrRecipeNotDeleted))
	assert.Equal(t, http.StatusInternalServerError, dbErrorStatus(errors.New("boom")))
}

func TestCreateRecipeForPerson(t *testing.T) {
	const body = `{"recipe": {"name": "Pie"}}`

	t.Run("a person can create their own recipe without an admin lookup", func(t *testing.T) {
		srv := recipeServer(nil, map[string]any{"recipe": "r"}, nil)
		c, w := jsonRequest(t, http.MethodPost, body)

		srv.CreateRecipeForPerson(c, 1, api.CreateRecipeForPersonParams{XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("an admin can create a recipe for the person they manage", func(t *testing.T) {
		srv := recipeServer([]error{nil}, map[string]any{"recipe": "r"}, nil)
		c, w := jsonRequest(t, http.MethodPost, body)

		srv.CreateRecipeForPerson(c, 2, api.CreateRecipeForPersonParams{XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("a family member who is not an admin cannot create a recipe for another person", func(t *testing.T) {
		srv := recipeServer([]error{errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodPost, body)

		srv.CreateRecipeForPerson(c, 2, api.CreateRecipeForPersonParams{XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCreateRecipeRelationship(t *testing.T) {
	t.Run("defaults the person to the caller", func(t *testing.T) {
		srv := recipeServer([]error{nil}, map[string]any{"Id": 1}, nil)
		c, w := jsonRequest(t, http.MethodPost, `{}`)

		srv.CreateRecipeRelationship(c, 5, api.CreateRecipeRelationshipParams{XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("the recipe must be visible to the caller", func(t *testing.T) {
		srv := recipeServer([]error{errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodPost, `{}`)

		srv.CreateRecipeRelationship(c, 5, api.CreateRecipeRelationshipParams{XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "recipe")
	})

	t.Run("an admin can like on behalf of the person they manage", func(t *testing.T) {
		srv := recipeServer([]error{nil, nil}, map[string]any{"Id": 1}, nil)
		c, w := jsonRequest(t, http.MethodPost, `{"person_id": 2}`)

		srv.CreateRecipeRelationship(c, 5, api.CreateRecipeRelationshipParams{XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("a non admin cannot like on behalf of another person", func(t *testing.T) {
		srv := recipeServer([]error{nil, errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodPost, `{"person_id": 2}`)

		srv.CreateRecipeRelationship(c, 5, api.CreateRecipeRelationshipParams{XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "person")
	})
}

func TestDeleteRecipeRelationship(t *testing.T) {
	t.Run("the caller removes their own like", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, nil)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.DeleteRecipeRelationship(c, 5, api.DeleteRecipeRelationshipParams{PersonId: 1, XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("the recipe must be visible to the caller", func(t *testing.T) {
		srv := recipeServer([]error{errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.DeleteRecipeRelationship(c, 5, api.DeleteRecipeRelationshipParams{PersonId: 1, XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("a non admin cannot remove another person's like", func(t *testing.T) {
		srv := recipeServer([]error{nil, errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.DeleteRecipeRelationship(c, 5, api.DeleteRecipeRelationshipParams{PersonId: 2, XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRecipeCommentNotFound(t *testing.T) {
	t.Run("update of a missing comment is 404", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, memgraph.ErrNotFound)
		c, w := jsonRequest(t, http.MethodPatch, `{"message": "hi"}`)

		srv.UpdateRecipeComment(c, 5, api.UpdateRecipeCommentParams{XUserID: 1})

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("delete of a missing comment is 404", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, memgraph.ErrNotFound)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.DeleteRecipeComment(c, 5, api.DeleteRecipeCommentParams{XUserID: 1})

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("other failures stay 500", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, errors.New("db error"))
		c, w := jsonRequest(t, http.MethodPatch, `{"message": "hi"}`)

		srv.UpdateRecipeComment(c, 5, api.UpdateRecipeCommentParams{XUserID: 1})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRecipeDeletion(t *testing.T) {
	t.Run("soft delete of an already deleted recipe is 404", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, memgraph.ErrNotFound)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.SoftDeleteRecipe(c, 5, api.SoftDeleteRecipeParams{XUserID: 1})

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("hard delete of a live recipe is 409", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, memgraph.ErrRecipeNotDeleted)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.HardDeleteRecipe(c, 5, api.HardDeleteRecipeParams{XUserID: 1})

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("hard delete of a soft deleted recipe succeeds", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, nil)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.HardDeleteRecipe(c, 5, api.HardDeleteRecipeParams{XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Recipe hard deleted")
	})

	t.Run("only managers may delete", func(t *testing.T) {
		srv := recipeServer([]error{errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodDelete, "")

		srv.HardDeleteRecipe(c, 5, api.HardDeleteRecipeParams{XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGetRecipe(t *testing.T) {
	t.Run("returns the recipe", func(t *testing.T) {
		srv := recipeServer([]error{nil, nil}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetRecipe(c, 5, api.GetRecipeParams{XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "read")
	})

	t.Run("requires visibility", func(t *testing.T) {
		srv := recipeServer([]error{errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetRecipe(c, 5, api.GetRecipeParams{XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing recipe is 404", func(t *testing.T) {
		srv := recipeServer([]error{nil, memgraph.ErrNotFound}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetRecipe(c, 5, api.GetRecipeParams{XUserID: 1})

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetFamilyCookbookDistance(t *testing.T) {
	for _, distance := range []int{0, -1, memgraph.MaxCookbookDistance + 1} {
		t.Run(fmt.Sprint(distance), func(t *testing.T) {
			srv := recipeServer(nil, nil, nil)
			c, w := jsonRequest(t, http.MethodGet, "")

			srv.GetFamilyCookbook(c, api.GetFamilyCookbookParams{Distance: distance, XUserID: 1})

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}

	t.Run("max distance is accepted", func(t *testing.T) {
		srv := recipeServer([]error{nil}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetFamilyCookbook(c, api.GetFamilyCookbookParams{Distance: memgraph.MaxCookbookDistance, XUserID: 1})

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetRecipesByPersonId(t *testing.T) {
	t.Run("requires visibility of the person", func(t *testing.T) {
		srv := recipeServer([]error{errDenied, errDenied}, nil, nil)
		c, w := jsonRequest(t, http.MethodGet, "")

		srv.GetRecipesByPersonId(c, 2, api.GetRecipesByPersonIdParams{XUserID: 1})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
