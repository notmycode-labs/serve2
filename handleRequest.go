package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join(*folder, r.URL.Path[1:])
	// Check if the requested path is within the specified folder
	relPath, err := filepath.Rel(*folder, filePath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		http.NotFound(w, r)
		return
	}

	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		renderDirectory(w, filePath)
		return
	}

	http.ServeFile(w, r, filePath)
}
