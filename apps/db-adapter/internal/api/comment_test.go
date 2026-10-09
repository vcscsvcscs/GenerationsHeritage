package api

import (
	"bytes"
	"context"
	"errors"
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

func mockServer(writeErr, authErr error, writeResult any) *server {
	mockSession := new(memgraphMock.SessionWithContext)
	mockDriver := new(memgraphMock.DriverWithContext)
	mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)

	if authErr != nil {
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, authErr)
	} else {
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	}

	if writeErr != nil {
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(nil, writeErr)
	} else {
		mockSession.On("ExecuteWrite", mock.Anything, mock.Anything, mock.Anything).Return(writeResult, nil)
	}

	mockSession.On("Close", mock.Anything).Return(nil)

	return &server{
		db:          mockDriver,
		dbOpTimeout: 2 * time.Second,
	}
}

func requestWithBody(method, url, body string) (*gin.Context, *httptest.ResponseRecorder) { //nolint:unparam // could be fixed in the future
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequestWithContext(context.Background(), method, url, bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func TestCommentOnPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		srv := mockServer(nil, nil, map[string]any{"ok": true})
		c, w := requestWithBody(http.MethodPost, "/comment", `{"text": "Hello"}`)
		params := api.CommentOnPersonParams{XUserID: 1}

		srv.CommentOnPerson(c, 123, params)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ok")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		srv := mockServer(nil, errors.New("unauthorized"), nil)
		c, w := requestWithBody(http.MethodPost, "/comment", `{"text": "Hi"}`)
		params := api.CommentOnPersonParams{XUserID: 1}

		srv.CommentOnPerson(c, 456, params)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("Database error", func(t *testing.T) {
		srv := mockServer(errors.New("db error"), nil, nil)
		c, w := requestWithBody(http.MethodPost, "/comment", `{"text": "Oops"}`)
		params := api.CommentOnPersonParams{XUserID: 2}

		srv.CommentOnPerson(c, 789, params)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "db error")
	})
}

func TestEditComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		srv := mockServer(nil, nil, map[string]any{"updated": true})
		c, w := requestWithBody(http.MethodPut, "/comment", `{"text": "Updated text"}`)
		params := api.EditCommentParams{XUserID: 4}

		srv.EditComment(c, 101, params)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "updated")
	})
}

func TestDeleteCommentOnPerson(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv := mockServer(nil, nil, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/comment", http.NoBody)

		params := api.DeleteCommentOnPersonParams{XUserID: 5}
		srv.DeleteCommentOnPerson(c, 202, params)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Comment deleted")
	})

	t.Run("DB error", func(t *testing.T) {
		srv := mockServer(errors.New("delete error"), nil, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/comment", http.NoBody)

		params := api.DeleteCommentOnPersonParams{XUserID: 6}
		srv.DeleteCommentOnPerson(c, 303, params)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "delete error")
	})
}

func TestGetCommentsOnPerson(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSession := new(memgraphMock.SessionWithContext)
		mockDriver := new(memgraphMock.DriverWithContext)
		mockDriver.On("NewSession", mock.Anything, mock.Anything).Return(mockSession)
		mockSession.On("ExecuteRead", mock.Anything, mock.Anything, mock.Anything).Return([]string{"Comment 1"}, nil)
		mockSession.On("Close", mock.Anything).Return(nil)

		srv := &server{
			db:          mockDriver,
			dbOpTimeout: 2 * time.Second,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/comment", http.NoBody)

		params := api.GetCommentsOnPersonParams{XUserID: 7}
		srv.GetCommentsOnPerson(c, 404, params)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Comment 1")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		srv := mockServer(nil, errors.New("access denied"), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/comment", http.NoBody)

		params := api.GetCommentsOnPersonParams{XUserID: 8}
		srv.GetCommentsOnPerson(c, 505, params)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "access denied")
	})
}
