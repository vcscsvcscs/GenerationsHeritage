package api

import (
	"errors"
	"fmt"
	"io"
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

func TestCreatePerson(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockResult := new(memgraphMock.Result)
		mockResult.On("Single", mock.Anything).Return(map[string]any{"id": 1}, nil)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("BeginTransaction", mock.Anything).Return(mockSession, nil)
		mockSession.On("Run", mock.Anything, memgraph.CreatePersonCypherQuery, mock.Anything).Return(mockResult, nil)
		mockSession.On("Run", mock.Anything, memgraph.CreateAdminRelationshipCypherQuery, mock.Anything).Return(nil, nil)
		mockSession.On("Commit", mock.Anything).Return(nil)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		body := `{"first_name": "test-person"}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(
			t.Context(), http.MethodPost, "/person", io.NopCloser(strings.NewReader(body)),
		)
		params := api.CreatePersonParams{XUserID: 1, XUserName: "test"}

		srv.CreatePerson(c, params)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "id")
	})

	t.Run("Bad request case", func(t *testing.T) {
		srv := &server{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/person", http.NoBody)

		srv.CreatePerson(c, api.CreatePersonParams{})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "msg")
	})

	t.Run("Internal server error case during transaction", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("BeginTransaction", mock.Anything).Return(nil, errors.New("transaction error"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		body := `{"first_name": "test-person"}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(
			t.Context(), http.MethodPost, "/person", io.NopCloser(strings.NewReader(body)),
		)
		params := api.CreatePersonParams{XUserID: 1}

		srv.CreatePerson(c, params)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "transaction error")
	})
}

func TestGetPersonById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything).Return(map[string]any{"id": 1}, nil)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/person/1", http.NoBody)
		params := api.GetPersonByIdParams{XUserID: 1}

		srv.GetPersonById(c, 1, params)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "id")
	})

	t.Run("Unauthorized case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unauthorized"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/person/1", http.NoBody)
		params := api.GetPersonByIdParams{XUserID: 2}

		srv.GetPersonById(c, 1, params)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

func TestSoftDeletePerson(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything).Return(nil, nil)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/person/1", http.NoBody)
		params := api.SoftDeletePersonParams{XUserID: 1}

		srv.SoftDeletePerson(c, 1, params)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "deleted")
	})

	t.Run("Unauthorized case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unauthorized"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/person/1", http.NoBody)
		params := api.SoftDeletePersonParams{XUserID: 2}

		srv.SoftDeletePerson(c, 1, params)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}
