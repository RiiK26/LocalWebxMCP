package mcp

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/RiiK26/LocalWebxMCP/Backend/internal/config"
)

// ErrPathEscape is returned when a requested filename would resolve
// outside the storage directory (e.g. via "../" traversal or an
// absolute path).
var ErrPathEscape = errors.New("filename escapes the storage directory")

// safeStoragePath resolves filename against the storage directory and
// guarantees the result stays inside it. Every tool that touches the
// filesystem based on AI-supplied input MUST go through this instead
// of calling filepath.Join(config.StorageDir, filename) directly.
func safeStoragePath(filename string) (string, error) {
	storageAbs, err := filepath.Abs(config.StorageDir)
	if err != nil {
		return "", err
	}

	// filepath.Join already calls Clean, which collapses "..",
	// but we resolve to an absolute path so a crafted absolute
	// filename (e.g. "/etc/passwd") can't bypass the check either.
	joined := filepath.Join(storageAbs, filename)
	resultAbs, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(storageAbs, resultAbs)
	if err != nil {
		return "", err
	}
	// If the relative path from storageAbs to resultAbs starts with
	// "..", resultAbs falls outside storageAbs.
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", ErrPathEscape
	}

	return resultAbs, nil
}
