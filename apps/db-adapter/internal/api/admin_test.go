package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	memgraphMock "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func TestCreateAdminRelationship(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(map[string]any{"result": "success"}, nil)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin", nil)
		params := api.CreateAdminRelationshipParams{XUserID: *api.IntPtr(1)}

		srv.CreateAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Unauthorized case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unauthorized"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin", nil)
		params := api.CreateAdminRelationshipParams{XUserID: *api.IntPtr(3)}

		srv.CreateAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("Internal server error case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin", nil)
		params := api.CreateAdminRelationshipParams{XUserID: 1}

		srv.CreateAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "db error")
	})
}

func TestDeleteAdminRelationship(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/admin", nil)
		params := api.DeleteAdminRelationshipParams{XUserID: 2}

		srv.DeleteAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "admin relationship was deleted")
	})

	t.Run("Unauthorized case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unauthorized"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/admin", nil)
		params := api.DeleteAdminRelationshipParams{XUserID: 3}

		srv.DeleteAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("Internal server error case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/admin", nil)
		params := api.DeleteAdminRelationshipParams{XUserID: 2}

		srv.DeleteAdminRelationship(c, 1, 2, params)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "db error")
	})
}
