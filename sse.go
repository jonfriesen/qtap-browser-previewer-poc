package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/qpoint-io/qreview/templates"
)

type SSEClient struct {
	w      http.ResponseWriter
	notify chan []Request
	done   chan bool
}

type SSEHub struct {
	clients map[*SSEClient]bool
	mutex   sync.RWMutex
}

var sseHub = &SSEHub{
	clients: make(map[*SSEClient]bool),
}

// Register a new SSE client
func (h *SSEHub) register(client *SSEClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.clients[client] = true
}

// Unregister an SSE client
func (h *SSEHub) unregister(client *SSEClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.clients, client)
	close(client.done)
}

// Broadcast request list updates to all connected clients
func (h *SSEHub) broadcast(requests []Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.notify <- requests:
		default:
			// Client channel is full, skip this client
		}
	}
}

// SSE endpoint handler
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create new client
	client := &SSEClient{
		w:      w,
		notify: make(chan []Request, 10),
		done:   make(chan bool),
	}

	// Register client
	sseHub.register(client)
	defer sseHub.unregister(client)

	// Send initial request list
	requestsMu.RLock()
	initialRequests := make([]Request, len(requests))
	copy(initialRequests, requests)
	requestsMu.RUnlock()

	sendRequestUpdate(client.w, initialRequests)

	// Handle client disconnect
	notify := r.Context().Done()

	// Listen for updates and disconnections
	for {
		select {
		case requestList := <-client.notify:
			if !sendRequestUpdate(client.w, requestList) {
				return
			}
		case <-notify:
			return
		case <-client.done:
			return
		}
	}
}

// Send request list update via SSE
func sendRequestUpdate(w http.ResponseWriter, requests []Request) bool {
	// Render the request list HTML
	var html strings.Builder
	templates.RequestList(convertToTemplateRequests(requests)).Render(context.Background(), &html)

	// Send as SSE event
	_, err := fmt.Fprintf(w, "event: request-update\ndata: %s\n\n", html.String())
	if err != nil {
		return false
	}

	// Flush the response
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return true
}

// Broadcast changes to all SSE clients
func broadcastRequestChanges() {
	requestsMu.RLock()
	defer requestsMu.RUnlock()

	// Create a copy of requests to send
	requestsCopy := make([]Request, len(requests))
	copy(requestsCopy, requests)

	sseHub.broadcast(requestsCopy)
}
