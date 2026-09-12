package mcp

import (
	"fmt"
	"os"
	"strings"

	"github.com/RiiK26/LocalWebxMCP/Backend/internal/config"
)

func schemaFindFiles() map[string]interface{} {
	return map[string]interface{}{
		"name":        "find_files",
		"description": "Searches for files by a partial name match (like 'find').",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{"type": "string", "description": "Part of the filename to search for"},
			},
			"required": []string{"query"},
		},
	}
}

func executeFindFiles(req *Request, resp *Response) {
	query, ok := req.Params.Arguments["query"].(string)
	if !ok {
		resp.Error = map[string]interface{}{"code": -32602, "message": "Invalid query parameter"}
		return
	}

	files, _ := os.ReadDir(config.StorageDir)
	var matchedFiles []string
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(query)) {
			matchedFiles = append(matchedFiles, f.Name())
		}
	}

	resultText := fmt.Sprintf("Found %d matching files: %v", len(matchedFiles), matchedFiles)
	resp.Result = map[string]interface{}{
		"content": []map[string]interface{}{{"type": "text", "text": resultText}},
	}
}
