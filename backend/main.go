package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}
