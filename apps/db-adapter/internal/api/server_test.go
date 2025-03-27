package api

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockHealthCheck struct {
	mock.Mock
}

func (m *mockHealthCheck) HealthCheckHandler(c *gin.Context) {
	m.Called(c)
}

func (m *mockHealthCheck) SetStatus(status string) {
	m.Called(status)
}

func (m *mockHealthCheck) GetStatus() string {
	args := m.Called()

	return args.String(0)
}

func TestNewServer(t *testing.T) {
	logger := zap.NewNop()
	mockDriver, err := neo4j.NewDriverWithContext("bolt+ssc://memgraph:7687", nil)
	assert.NoError(t, err)
	mockHealth := &mockHealthCheck{}

	t.Run("should create a new server instance", func(t *testing.T) {
		t.Parallel()
		srv := New(logger, mockDriver, mockHealth, 1)
		assert.NotNil(t, srv)
	})

	t.Run("should panic if logger is nil", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			New(nil, mockDriver, mockHealth, 1)
		})
	})

	t.Run("should panic if driver is nil", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			New(logger, nil, mockHealth, 1)
		})
	})

	t.Run("should panic if healthcheck is nil", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			New(logger, mockDriver, nil, 1)
		})
	})

	t.Run("should panic if databaseOperationTimeout is 0", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			New(logger, mockDriver, mockHealth, 0)
		})
	})
}

func TestHealthCheck(t *testing.T) {
	mockHealth := &mockHealthCheck{}
	srv := &server{health: mockHealth}

	t.Run("should call HealthCheckHandler", func(t *testing.T) {
		mockHealth.On("HealthCheckHandler", mock.Anything).Once()

		c, _ := gin.CreateTestContext(nil)
		srv.HealthCheck(c)

		mockHealth.AssertCalled(t, "HealthCheckHandler", c)
	})
}
