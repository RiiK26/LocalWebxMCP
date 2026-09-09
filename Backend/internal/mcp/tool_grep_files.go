package mcp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func schemaGrepFile() map[string]interface{} {
	return map[string]interface{}{
		"name":        "grep_file",
		"description": "Searches for specific text inside a file and returns matching lines (like 'grep').",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"filename": map[string]interface{}{"type": "string", "description": "The file to search in"},
				"pattern":  map[string]interface{}{"type": "string", "description": "The text pattern to search for"},
			},
			"required": []string{"filename", "pattern"},
		},
	}
}

func executeGrepFile(req *Request, resp *Response) {
	filename, ok1 := req.Params.Arguments["filename"].(string)
	pattern, ok2 := req.Params.Arguments["pattern"].(string)

	if !ok1 || !ok2 {
		resp.Error = map[string]interface{}{"code": -32602, "message": "Invalid parameters for grep"}
		return
	}

	safePath, err := safeStoragePath(filename)
	if err != nil {
		resp.Error = map[string]interface{}{"code": -32602, "message": "Invalid filename: " + err.Error()}
		return
	}

	file, err := os.Open(safePath)
	if err != nil {
		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": "Error: File not found"}},
		}
		return
	}
	defer file.Close()

	var matchedLines []string
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			matchedLines = append(matchedLines, fmt.Sprintf("Line %d: %s", lineNumber, line))
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		matchedLines = append(matchedLines, fmt.Sprintf("Error reading file midway: %v", err))
	}

	resultText := strings.Join(matchedLines, "\n")
	if len(resultText) == 0 {
		resultText = "No matches found for the given pattern."
	} else if len(resultText) > 50000 {
		resultText = resultText[:50000] + "\n...(truncated because too large)"
	}

	resp.Result = map[string]interface{}{
		"content": []map[string]interface{}{{"type": "text", "text": resultText}},
	}
}
