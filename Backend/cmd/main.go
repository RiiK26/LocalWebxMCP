package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ItsMe-RiiK/LocalWebxMCP/Backend/internal/api"
	"github.com/ItsMe-RiiK/LocalWebxMCP/Backend/internal/config"
	"github.com/ItsMe-RiiK/LocalWebxMCP/Backend/internal/mcp"
	"github.com/ItsMe-RiiK/LocalWebxMCP/Backend/internal/middleware"
	"github.com/ItsMe-RiiK/LocalWebxMCP/Backend/internal/storage"
)

func main() {
	// 1. Initialize Storage Directory
	if err := storage.Init(); err != nil {
		log.Fatalf("Failed to initialize storage directory: %v", err)
	}

	// 2. Create Router
	mux := http.NewServeMux()

	// --- FRONTEND & API ROUTES ---
	mux.Handle("/", http.FileServer(http.Dir(config.FrontendDir)))
	mux.Handle("/api/files/", http.StripPrefix("/api/files/", http.FileServer(http.Dir(config.StorageDir))))
	mux.HandleFunc("/api/upload", api.UploadHandler)
	mux.HandleFunc("/api/files", api.ListFilesHandler)

	// --- AI MCP ROUTE ---
	mux.HandleFunc("/mcp", mcp.Handler)
	mux.HandleFunc("/mcp/", mcp.Handler)

	// 3. Apply Security Middleware
	secureMux := middleware.CORS(mux)

	// 4. Start Server
	fmt.Printf("UI Server running at http://localhost%s\n", config.ServerPort)
	fmt.Printf("MCP Endpoint active at http://localhost%s/mcp\n", config.ServerPort)

	err := http.ListenAndServe(config.ServerPort, secureMux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
