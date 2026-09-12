package storage

import (
	"os"

	"github.com/RiiK26/LocalWebxMCP/Backend/internal/config"
)

// Init creates the necessary storage directory if it doesn't exist
func Init() error {
	return os.MkdirAll(config.StorageDir, os.ModePerm)
}
