package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	apiServer "github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/api"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/internal/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/gin/healthcheck"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	defaultHTTPPort       = ":80"
	defaultMemgraphURI    = "bolt://memgraph:7687"
	defaultMemgraphUser   = "memgraph"
	defaultMemgraphPass   = "memgraph"
	defaultProduction     = false
	defaultRequestTimeout = 20
	defaultDBOpTimeout    = 5
)

var (
	httpPort       string
	memgraphURI    string
	memgraphUser   string
	memgraphPass   string
	production     bool
	requestTimeout time.Duration
	dbOpTimeout    time.Duration
)

func init() { //nolint:gochecknoinits // this is a main package, init is ok
	viper.AutomaticEnv()

	viper.SetDefault("HTTP_PORT", defaultHTTPPort)
	viper.SetDefault("MEMGRAPH_URI", defaultMemgraphURI)
	viper.SetDefault("MEMGRAPH_USER", defaultMemgraphUser)
	viper.SetDefault("MEMGRAPH_PASS", defaultMemgraphPass)
	viper.SetDefault("PRODUCTION", defaultProduction)
	viper.SetDefault("REQUEST_TIMEOUT", defaultRequestTimeout)
	viper.SetDefault("DB_OP_TIMEOUT", defaultDBOpTimeout)

	httpPort = viper.GetString("HTTP_PORT")
	memgraphURI = viper.GetString("MEMGRAPH_URI")
	memgraphUser = viper.GetString("MEMGRAPH_USER")
	memgraphPass = viper.GetString("MEMGRAPH_PASS")
	production = viper.GetBool("PRODUCTION")
	requestTimeout = time.Duration(viper.GetInt("REQUEST_TIMEOUT")) * time.Second
	dbOpTimeout = time.Duration(viper.GetInt("DB_OP_TIMEOUT")) * time.Millisecond
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

	memgraphDriver := memgraph.InitDatabase(logger, memgraphURI, memgraphUser, memgraphPass)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://localhost%s", httpPort), "http://localhost"},
		AllowCredentials: true,
		AllowHeaders:     []string{"X-User-ID", "X-User-Name", "Content-Type"},
		MaxAge:           12 * time.Hour,
	}))

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	sApi := apiServer.New(logger, memgraphDriver, hc, dbOpTimeout)
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
