package mcp

import (
	"fmt"
	"os"
	"strings"

	"github.com/RiiK26/LocalWebxMCP/Backend/internal/config"
)

func schemaListFiles() map[string]interface{} {
	return map[string]interface{}{
		"name":        "list_files",
		"description": "Lists all files in the local storage including metadata like size and last modified date (like 'ls -la').",
		"inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
	}
}

func executeListFiles(_ *Request, resp *Response) {
	files, err := os.ReadDir(config.StorageDir)
	if err != nil {
		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": "Error: Cannot read directory"}},
		}
		return
	}

	var fileDetails []string
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			continue
		}

		size := info.Size()
		sizeStr := fmt.Sprintf("%d bytes", size)
		if size > 1024*1024 {
			sizeStr = fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
		} else if size > 1024 {
			sizeStr = fmt.Sprintf("%.2f KB", float64(size)/1024)
		}

		modTime := info.ModTime().Format("2006-01-02 15:04:05")
		fileDetails = append(fileDetails, fmt.Sprintf("- %s | Size: %s | Modified: %s", f.Name(), sizeStr, modTime))
	}

	resultText := "Storage is empty."
	if len(fileDetails) > 0 {
		resultText = "Files in storage:\n" + strings.Join(fileDetails, "\n")
	}

	resp.Result = map[string]interface{}{
		"content": []map[string]interface{}{{"type": "text", "text": resultText}},
	}
}
