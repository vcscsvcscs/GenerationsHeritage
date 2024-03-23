package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/backend/liveness"
)

var (
	cert            = flag.String("cert", "./private/keys/cert.pem", "Specify the path of TLS cert")
	key             = flag.String("key", "./private/keys/key.pem", "Specify the path of TLS key")
	httpsPort       = flag.String("https", ":443", "Specify port for http secure hosting(example for format :443)")
	httpPort        = flag.String("http", ":80", "Specify port for http hosting(example for format :80)")
	release         = flag.Bool("release", false, "Set true to release build")
	logToFile       = flag.Bool("log-to-file", false, "Set true to log to file")
	logToFileAndStd = flag.Bool("log-to-file-and-std", false, "Set true to log to file and std")
)

func main() {
	flag.Parse()
	if *release {
		gin.SetMode(gin.ReleaseMode)
	}
	if *logToFileAndStd || *logToFile {
		gin.DisableConsoleColor() // Disable Console Color, you don't need console color when writing the logs to file.
		path := fmt.Sprintf("private/logs/%02dy_%02dm_%02dd_%02dh_%02dm_%02ds.log", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		//fmt.Println(path)
		logerror1 := os.MkdirAll("private/logs/", 0755)
		f, logerror2 := os.Create(path)
		if logerror1 != nil || logerror2 != nil {
			log.Println("Cant log to file")
		} else if *logToFileAndStd {
			gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		} else {
			gin.DefaultWriter = io.MultiWriter(f)
		}
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(gin.DefaultErrorWriter)

	hc := liveness.New()

	router := gin.Default()
	router.GET("/health", hc.HealthCheckHandler())

	//Server configuration
	var server *http.Server

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
