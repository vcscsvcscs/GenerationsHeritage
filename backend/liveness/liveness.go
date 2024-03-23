package liveness

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type HealthCheck interface {
	SetStatus(status string)
	GetStatus() string
	HealthCheckHandler() gin.HandlerFunc
}

type healthCheck struct {
	status string
	sync   sync.Mutex
}

func New() HealthCheck {
	return &healthCheck{
		status: "ok",
	}
}

func (hc *healthCheck) SetStatus(status string) {
	hc.sync.Lock()
	defer hc.sync.Unlock()

	hc.status = status
}

func (hc *healthCheck) GetStatus() string {
	hc.sync.Lock()
	defer hc.sync.Unlock()

	return hc.status
}

func (hc *healthCheck) HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch hc.GetStatus() {
		case "nok":
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": hc.GetStatus(),
			})
		default:
			c.JSON(http.StatusOK, gin.H{
				"status": hc.GetStatus(),
			})
		}
	}
}
