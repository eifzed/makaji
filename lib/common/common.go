package common

import (
	"crypto/rand"
	"io"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/oklog/ulid"
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

func GenerateUUIDV7() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)

	ulid := ulid.MustNew(ulid.Timestamp(t), entropy)
	return ulid.String()
}
