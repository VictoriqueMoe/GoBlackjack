package utils

import (
	"fmt"
	"os"
)

// ConnectionURL func for building URL connection.
func ConnectionURL() string {
	// Return connection URL.
	return fmt.Sprintf(
		"%s:%s",
		os.Getenv("SERVER_HOST"),
		os.Getenv("SERVER_PORT"),
	)
}
