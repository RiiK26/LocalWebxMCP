package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/RiiK26/LocalWebxMCP/Backend/internal/config"
)

// UploadHandler handles file uploads from the frontend
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB limit
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join(config.StorageDir, handler.Filename))
	if err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)
	w.Write([]byte(`{"status":"success"}`))
}

// ListFilesHandler returns a JSON array of all uploaded files
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := os.ReadDir(config.StorageDir)
	var fileNames []string

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"files": fileNames})
}
