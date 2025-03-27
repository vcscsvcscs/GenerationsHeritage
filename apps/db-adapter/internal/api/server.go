package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/gin/healthcheck"
	"go.uber.org/zap"
)

type server struct {
	db          neo4j.DriverWithContext
	dbOpTimeout time.Duration
	health      healthcheck.HealthCheck
	logger      *zap.Logger
}

func New(logger *zap.Logger, drv neo4j.DriverWithContext, healthcheck healthcheck.HealthCheck, databaseOperationTimeoutInMs time.Duration) ServerInterface {
	if logger == nil {
		panic("logger is required")
	}

	if drv == nil {
		panic("neo4j driver is required")
	}

	if healthcheck == nil {
		panic("healthcheck is required")
	}

	if databaseOperationTimeoutInMs == 0 {
		panic("database operation timeout is required")
	}

	return &server{db: drv, health: healthcheck, logger: logger}
}

func (srv *server) HealthCheck(c *gin.Context) {
	srv.health.HealthCheckHandler(c)
}
