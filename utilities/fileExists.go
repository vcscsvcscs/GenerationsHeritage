package utilities

import (
	"os"
)

// Checks if file on path exists or not
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
