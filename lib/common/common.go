package common

import (
	"io"
	"os"

	"github.com/labstack/gommon/log"
)

func SafelyCloseFile(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Warnf("Failed to close file: %s\n", err)
	}
}

func IsDevelopment() bool {
	isLocal := os.Getenv("ISLOCAL")
	return isLocal == "1"
}
