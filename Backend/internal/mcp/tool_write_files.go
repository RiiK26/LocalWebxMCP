package mcp

import (
	"os"
)

func schemaWriteFile() map[string]interface{} {
	return map[string]interface{}{
		"name":        "write_file",
		"description": "Creates a new file or overwrites an existing file with the provided text content.",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"filename": map[string]interface{}{
					"type":        "string",
					"description": "The name of the file to create/overwrite (e.g., script.js)",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "The exact text content to write into the file",
				},
			},
			"required": []string{"filename", "content"},
		},
	}
}

func executeWriteFile(req *Request, resp *Response) {
	filename, ok1 := req.Params.Arguments["filename"].(string)
	content, ok2 := req.Params.Arguments["content"].(string)

	if !ok1 || !ok2 {
		resp.Error = map[string]interface{}{"code": -32602, "message": "Invalid parameters for write_file"}
		return
	}

	filePath, err := safeStoragePath(filename)
	if err != nil {
		resp.Error = map[string]interface{}{"code": -32602, "message": "Invalid filename: " + err.Error()}
		return
	}

	// 0644 means: Read/Write for owner, Read-only for others
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": "Error: Failed to write to file"}},
		}
		return
	}

	resp.Result = map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": "Success: File '" + filename + "' has been saved successfully."},
		},
	}
}
