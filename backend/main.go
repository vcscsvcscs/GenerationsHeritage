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
	"github.com/vcscsvcscs/GenerationsHeritage/backend/handlers"
	"github.com/vcscsvcscs/GenerationsHeritage/backend/memgraph"
	"github.com/vcscsvcscs/GenerationsHeritage/utilities"
	"github.com/vcscsvcscs/GenerationsHeritage/utilities/gin_liveness"
)

var (
	cert            = flag.String("cert", "./private/keys/cert.pem", "Specify the path of TLS cert")
	key             = flag.String("key", "./private/keys/key.pem", "Specify the path of TLS key")
	httpsPort       = flag.String("https", ":443", "Specify port for http secure hosting(example for format :443)")
	httpPort        = flag.String("http", ":80", "Specify port for http hosting(example for format :80)")
	memgraphURI     = flag.String("memgraph", "bolt+ssc://memgraph:7687", "Specify the Memgraph database URI")
	memgraphUser    = flag.String("memgraph-user", "", "Specify the Memgraph database user")
	memgraphPass    = flag.String("memgraph-pass", "", "Specify the Memgraph database password")
	release         = flag.Bool("release", false, "Set true to release build")
	logToFile       = flag.Bool("log-to-file", false, "Set true to log to file")
	logToFileAndStd = flag.Bool("log-to-file-and-std", false, "Set true to log to file and std")
	requestTimeout  = time.Duration(*flag.Int("request-timeout", 20, "Set request timeout in seconds"))
)

func main() {
	flag.Parse()
	if *release {
		gin.SetMode(gin.ReleaseMode)
	}

	utilities.SetupLogger(*logToFileAndStd, *logToFile)

	hc := gin_liveness.New()

	memgraphDriver := memgraph.InitDatabase(*memgraphURI, *memgraphUser, *memgraphPass)

	router := gin.Default()
	router.GET("/health", hc.HealthCheckHandler())
	router.GET("/person", handlers.ViewPerson(memgraphDriver))
	router.POST("/createPerson", handlers.CreatePerson(memgraphDriver))
	router.DELETE("/deletePerson", handlers.DeletePerson(memgraphDriver))
	router.PUT("/updatePerson", handlers.UpdatePerson(memgraphDriver))
	router.POST("/createRelationship", handlers.CreateRelationship(memgraphDriver))
	router.DELETE("/deleteRelationship", handlers.DeleteRelationship(memgraphDriver))
	router.PUT("/verifyRelationship", handlers.VerifyRelationship(memgraphDriver))
	router.POST("/createRelationshipAndPerson", handlers.CreateRelationshipAndPerson(memgraphDriver))

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
