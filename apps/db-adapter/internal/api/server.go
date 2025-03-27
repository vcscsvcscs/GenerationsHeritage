package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/gin/healthcheck"
	"go.uber.org/zap"
)

type server struct {
	db     neo4j.DriverWithContext
	health healthcheck.HealthCheck
	logger *zap.Logger
}

func New(logger *zap.Logger, drv neo4j.DriverWithContext, healthcheck healthcheck.HealthCheck) ServerInterface {
	if logger == nil {
		panic("logger is required")
	}

	if drv == nil {
		panic("neo4j driver is required")
	}

	if healthcheck == nil {
		panic("healthcheck is required")
	}

	return &server{db: drv, health: healthcheck, logger: logger}
}

func (srv *server) HealthCheck(c *gin.Context) {
	srv.health.HealthCheckHandler(c)
}
