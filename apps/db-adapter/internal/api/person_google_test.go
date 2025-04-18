package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	memgraphMock "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph/mock"
)

func TestGetPersonByGoogleId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(
			map[string]any{"person": "test-person"}, nil,
		)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(
			t.Context(), http.MethodGet, "/person/google-id", http.NoBody,
		)

		srv.GetPersonByGoogleId(c, "test-google-id")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test-person")
	})

	t.Run("Internal server error case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/person/google-id", http.NoBody)

		srv.GetPersonByGoogleId(c, "test-google-id")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "db error")
	})
}

func TestCreatePersonByGoogleIdAndInviteCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(map[string]any{"result": "success"}, nil)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"invite_code": "test-code", "person": {"first_name": "test-person"}}`
		c.Request = httptest.NewRequest(
			http.MethodPost, "/person/google-id/invite-code", io.NopCloser(strings.NewReader(body)),
		)

		srv.CreatePersonByGoogleIdAndInviteCode(c, "test-google-id")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Bad request case", func(t *testing.T) {
		srv := &server{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"invalid_json":`
		c.Request = httptest.NewRequest(
			http.MethodPost, "/person/google-id/invite-code", io.NopCloser(strings.NewReader(body)),
		)

		srv.CreatePersonByGoogleIdAndInviteCode(c, "test-google-id")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"msg\":\"unexpected EOF\"}")
	})
}

func TestCreatePersonByGoogleId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful case", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(
			map[string]any{"result": "success"}, nil,
		)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 5 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"name": "test-person"}`
		c.Request = httptest.NewRequestWithContext(
			t.Context(), http.MethodPost, "/person/google-id", io.NopCloser(strings.NewReader(body)),
		)

		srv.CreatePersonByGoogleId(c, "test-google-id")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Bad request case", func(t *testing.T) {
		srv := &server{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"invalid_json":`
		c.Request = httptest.NewRequestWithContext(
			t.Context(), http.MethodPost, "/person/google-id", io.NopCloser(strings.NewReader(body)),
		)

		srv.CreatePersonByGoogleId(c, "test-google-id")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"msg\":\"unexpected EOF\"}")
	})
}
