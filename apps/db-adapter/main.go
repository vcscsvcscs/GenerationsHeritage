package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/api"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/gin/healthcheck"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/memgraph"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	httpPort       string
	memgraphURI    string
	memgraphUser   string
	memgraphPass   string
	production     bool
	requestTimeout time.Duration
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("HTTP_PORT", ":80")
	viper.SetDefault("MEMGRAPH_URI", "bolt+ssc://memgraph:7687")
	viper.SetDefault("MEMGRAPH_USER", "")
	viper.SetDefault("MEMGRAPH_PASS", "")
	viper.SetDefault("PRODUCTION", false)
	viper.SetDefault("REQUEST_TIMEOUT", 20)
	viper.SetDefault("DB_OP_TIMEOUT", 5)

	httpPort = viper.GetString("HTTP_PORT")
	memgraphURI = viper.GetString("MEMGRAPH_URI")
	memgraphUser = viper.GetString("MEMGRAPH_USER")
	memgraphPass = viper.GetString("MEMGRAPH_PASS")
	production = viper.GetBool("PRODUCTION")
	requestTimeout = time.Duration(viper.GetInt("REQUEST_TIMEOUT")) * time.Second
}

func main() {
	var logger *zap.Logger
	var err error
	if production {
		logger, err = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	hc := healthcheck.New()

	memgraphDriver := memgraph.InitDatabase(memgraphURI, memgraphUser, memgraphPass)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost"},
		AllowCredentials: true,
		AllowHeaders:     []string{"X-User-ID", "Content-Type"},
		MaxAge:           12 * time.Hour,
	}))

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	sApi := api.New(logger, memgraphDriver, hc)
	api.RegisterHandlersWithOptions(router, sApi, api.GinServerOptions{})

	server := &http.Server{
		Addr:         httpPort,
		Handler:      router,
		ReadTimeout:  requestTimeout * time.Second,
		WriteTimeout: requestTimeout * time.Second,
	}
	go func() {
		logger.Info("Starting server", zap.String("port", httpPort))
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with some time to finish requests.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	logger.Info("Shutting down server...")

	// The context is used to inform the server it has some seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*time.Second)
	defer cancel()

	go func() {
		if err := memgraphDriver.Close(context.WithoutCancel(ctx)); err != nil {
			logger.Error("Memgraph driver forced to shutdown", zap.Error(err))
		}
	}()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
