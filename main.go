package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static/**
var staticFS embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files
	embeddedStatic, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("failed to prepare embedded static fs: %v", err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(embeddedStatic))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/requests/", requestHandler)
	http.HandleFunc("/select-request", selectRequestHandler)
	http.HandleFunc("/events", sseHandler)

	// Start request generator
	startRequestGenerator()

	fmt.Printf("Server starting on http://0.0.0.0:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
