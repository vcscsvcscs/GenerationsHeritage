package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	utilities "github.com/vcscsvcscs/GenerationsHeritage/pkg"
	"github.com/vcscsvcscs/GenerationsHeritage/pkg/gin/healthcheck"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

var (
	cert             = flag.String("cert", "/etc/gh-auth-service/ssl/tls.crt", "Specify the path of TLS cert")
	key              = flag.String("key", "/etc/gh-auth-service/ssl/tls.key", "Specify the path of TLS key")
	zitadelAccessKey = flag.String("zitadel-access-key", "/etc/gh-auth-service/zitadel/api-key.json", "Specify the path of Zitadel access key")
	httpsPort        = flag.String("https", ":443", "Specify port for http secure hosting(example for format :443)")
	httpPort         = flag.String("http", ":80", "Specify port for http hosting(example for format :80)")
	zitadelURI       = flag.String("zitadel-uri", "zitadel.varghacsongor.hu", "Specify the Zitadel URI")
	release          = flag.Bool("release", false, "Set true to release build")
	logToFile        = flag.Bool("log-to-file", false, "Set true to log to file")
	logToFileAndStd  = flag.Bool("log-to-file-and-std", false, "Set true to log to file and std")
	requestTimeout   = time.Duration(*flag.Int("request-timeout", 20, "Set request timeout in seconds"))
)

func main() {
	flag.Parse()
	if *release {
		gin.SetMode(gin.ReleaseMode)
	}

	utilities.SetupLogger(*logToFileAndStd, *logToFile)

	hc := healthcheck.New()

	router := gin.Default()
	router.Use(gin.Recovery())

	ctx := context.Background()

	// Initiate the authorization by providing a zitadel configuration and a verifier.
	// This example will use OAuth2 Introspection for this, therefore you will also need to provide the downloaded api key.json
	authZ, err := authorization.New(ctx, zitadel.New(*zitadelURI), oauth.DefaultAuthorization(*zitadelAccessKey))
	if err != nil {
		log.Println("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	// Initialize the HTTP middleware by providing the authorization
	mw := middleware.New(authZ)

	router.GET("/health", hc.HealthCheckHandler())
	router.GET("/auth", auth(mw))

	server := utilities.SetupHttpsServer(router, *cert, *key, *httpsPort, *httpPort, requestTimeout)

	// Wait for interrupt signal to gracefully shutdown the server with some time to finish requests.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has some seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
