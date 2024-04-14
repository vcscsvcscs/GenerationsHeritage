package utilities

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupLogger(logToFileAndStd bool, logToFile bool) {
	if logToFileAndStd || logToFile {
		gin.DisableConsoleColor() // Disable Console Color, you don't need console color when writing the logs to file.
		path := fmt.Sprintf("private/logs/%02dy_%02dm_%02dd_%02dh_%02dm_%02ds.log", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		logerror1 := os.MkdirAll("private/logs/", 0755)
		f, logerror2 := os.Create(path)
		if logerror1 != nil || logerror2 != nil {
			log.Println("Cant log to file")
		} else if logToFileAndStd {
			gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		} else {
			gin.DefaultWriter = io.MultiWriter(f)
		}
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(gin.DefaultErrorWriter)
}
