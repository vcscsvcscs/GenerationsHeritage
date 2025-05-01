package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/gin/healthcheck"
	"go.uber.org/zap"
)

type server struct {
	logger      *zap.Logger
	db          neo4j.DriverWithContext
	health      healthcheck.HealthCheck
	dbOpTimeout time.Duration
}

func New(
	logger *zap.Logger, drv neo4j.DriverWithContext, hc healthcheck.HealthCheck, databaseOperationTimeout time.Duration,
) api.ServerInterface {
	if logger == nil {
		panic("logger is required")
	}

	if drv == nil {
		panic("neo4j driver is required")
	}

	if hc == nil {
		panic("healthcheck is required")
	}

	if databaseOperationTimeout == 0 {
		panic("database operation timeout is required")
	}

	return &server{db: drv, health: hc, logger: logger, dbOpTimeout: databaseOperationTimeout}
}

func (srv *server) HealthCheck(c *gin.Context) {
	srv.health.HealthCheckHandler(c)
}

// Helper function to create a session with timeout
func (srv *server) createSessionWithTimeout(ctx context.Context) neo4j.SessionWithContext {
	actx, acancel := context.WithTimeout(ctx, srv.dbOpTimeout)
	defer acancel()
	session := srv.db.NewSession(actx, neo4j.SessionConfig{})
	return session
}
