package main

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/qpoint-io/qreview/templates"
)

type Request struct {
	ConnectionID          string   `json:"connectionId,omitzero"`
	EndpointId            string   `json:"endpointId,omitzero"`
	RequestId             string   `json:"requestId,omitzero"`
	TLSProbeTypesDetected []string `json:"tlsProbeTypesDetected,omitzero"`
	TLSIntrospected       bool     `json:"tlsProbeIntrospected,omitzero"`
	Tags                  []string `json:"tags,omitempty"`

	Timestamp   time.Time `json:"timestamp"`
	Direction   string    `json:"direction"`
	Url         string    `json:"url"`
	URLPath     string    `json:"path"`
	Method      string    `json:"method"`
	Status      int       `json:"status"`
	Duration    int64     `json:"duration"`
	ContentType string    `json:"contentType"`
	Category    string    `json:"category"`
	Agent       string    `json:"agent"`

	WrBytes int64 `json:"bytesSent"`
	RdBytes int64 `json:"bytesReceived"`

	AuthTokenMask string `json:"authTokenMask"`
	// AuthTokenHash is a SHA-256 hash of the auth token. The length is 32 bytes (64 characters) enforced by ClickHouse.
	AuthTokenHash   string `json:"authTokenHash"`
	AuthTokenSource string `json:"authTokenSource"`
	AuthTokenType   string `json:"authTokenType"`
}

var (
	requests   []Request
	requestsMu sync.RWMutex
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	requestsMu.RLock()
	defer requestsMu.RUnlock()

	templates.NetworkPage(convertToTemplateRequests(requests)).Render(r.Context(), w)
}

// clearRequestsHandler clears all requests (for testing)
func clearRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestsMu.Lock()
	requests = []Request{}
	requestsMu.Unlock()

	requestsMu.RLock()
	defer requestsMu.RUnlock()
	templates.RequestList(convertToTemplateRequests(requests)).Render(r.Context(), w)

	// Broadcast changes to SSE clients
	go broadcastRequestChanges()
}

// selectRequestHandler handles row click to show split view
func selectRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	r.ParseForm()
	url := r.FormValue("url")
	statusStr := r.FormValue("status")
	method := r.FormValue("method")
	contentType := r.FormValue("contentType")
	durationStr := r.FormValue("duration")
	rdBytesStr := r.FormValue("rdBytes")
	agent := r.FormValue("agent")

	// Convert strings to appropriate types
	status, _ := strconv.Atoi(statusStr)
	duration, _ := strconv.ParseInt(durationStr, 10, 64)
	rdBytes, _ := strconv.ParseInt(rdBytesStr, 10, 64)

	// Create selected request
	selectedRequest := &templates.Request{
		Url:         url,
		Method:      method,
		Status:      status,
		ContentType: contentType,
		Duration:    duration,
		RdBytes:     rdBytes,
		Agent:       agent,
		Timestamp:   time.Now(),
	}

	// Get current requests for the left panel
	requestsMu.RLock()
	templateRequests := convertToTemplateRequests(requests)
	requestsMu.RUnlock()

	// Render split view
	templates.NetworkMainContent(templateRequests, selectedRequest).Render(r.Context(), w)
}

// Route handler for matching request routes
func requestHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/requests/clear" && r.Method == http.MethodPost {
		clearRequestsHandler(w, r)
		return
	}

	http.NotFound(w, r)
}

// Request generator

// generateRandomRequest creates a realistic random egress API request
func generateRandomRequest() Request {
	// Real-world API endpoints and domains
	apiEndpoints := []struct {
		domain      string
		path        string
		contentType string
		status      int
		duration    int64
		bytes       int64
	}{
		// Payment APIs
		{"api.stripe.com", "/v1/payment_intents", "application/json", 200, 300 + rand.Int64N(500), 1024 + rand.Int64N(2048)},
		{"api.paypal.com", "/v2/checkout/orders", "application/json", 201, 400 + rand.Int64N(600), 2048 + rand.Int64N(4096)},
		{"api.squareup.com", "/v2/payments", "application/json", 200, 250 + rand.Int64N(400), 1536 + rand.Int64N(3072)},

		// Cloud Services
		{"api.aws.amazon.com", "/s3/2019-10-31/bucket", "application/xml", 200, 200 + rand.Int64N(300), 512 + rand.Int64N(1024)},
		{"storage.googleapis.com", "/upload/storage/v1/b", "application/json", 200, 150 + rand.Int64N(250), 256 + rand.Int64N(512)},
		{"api.cloudflare.com", "/client/v4/zones", "application/json", 200, 100 + rand.Int64N(200), 1024 + rand.Int64N(2048)},

		// Social Media APIs
		{"api.twitter.com", "/2/tweets", "application/json", 201, 500 + rand.Int64N(800), 2048 + rand.Int64N(4096)},
		{"graph.facebook.com", "/v18.0/me", "application/json", 200, 300 + rand.Int64N(500), 512 + rand.Int64N(1024)},
		{"api.linkedin.com", "/v2/ugcPosts", "application/json", 201, 400 + rand.Int64N(600), 1536 + rand.Int64N(3072)},

		// Communication APIs
		{"api.twilio.com", "/2010-04-01/Accounts", "application/x-www-form-urlencoded", 201, 600 + rand.Int64N(1000), 256 + rand.Int64N(512)},
		{"api.sendgrid.com", "/v3/mail/send", "application/json", 202, 200 + rand.Int64N(400), 128 + rand.Int64N(256)},
		{"api.slack.com", "/api/chat.postMessage", "application/json", 200, 150 + rand.Int64N(300), 512 + rand.Int64N(1024)},

		// Analytics & Monitoring
		{"www.google-analytics.com", "/collect", "text/plain", 200, 50 + rand.Int64N(100), 64 + rand.Int64N(128)},
		{"api.mixpanel.com", "/track", "application/json", 200, 100 + rand.Int64N(200), 128 + rand.Int64N(256)},
		{"api.segment.io", "/v1/track", "application/json", 200, 80 + rand.Int64N(150), 256 + rand.Int64N(512)},

		// Maps & Location
		{"maps.googleapis.com", "/maps/api/geocode/json", "application/json", 200, 200 + rand.Int64N(400), 1024 + rand.Int64N(2048)},
		{"api.mapbox.com", "/geocoding/v5/mapbox.places", "application/json", 200, 150 + rand.Int64N(300), 1536 + rand.Int64N(3072)},

		// AI/ML Services
		{"api.openai.com", "/v1/chat/completions", "application/json", 200, 2000 + rand.Int64N(5000), 4096 + rand.Int64N(8192)},
		{"api.anthropic.com", "/v1/messages", "application/json", 200, 1500 + rand.Int64N(4000), 3072 + rand.Int64N(6144)},
		{"api.huggingface.co", "/models", "application/json", 200, 1000 + rand.Int64N(3000), 2048 + rand.Int64N(4096)},

		// Database & Storage
		{"api.supabase.co", "/rest/v1", "application/json", 200, 100 + rand.Int64N(200), 512 + rand.Int64N(1024)},
		{"api.firebase.google.com", "/v1beta1/projects", "application/json", 200, 150 + rand.Int64N(300), 1024 + rand.Int64N(2048)},

		// CDN & Assets
		{"cdn.jsdelivr.net", "/npm", "application/javascript", 200, 50 + rand.Int64N(100), 4096 + rand.Int64N(8192)},
		{"unpkg.com", "/", "application/javascript", 200, 80 + rand.Int64N(150), 2048 + rand.Int64N(4096)},
	}

	// Error scenarios for external APIs
	errorEndpoints := []struct {
		domain      string
		path        string
		contentType string
		status      int
		duration    int64
	}{
		{"api.stripe.com", "/v1/payment_intents", "application/json", 402, 100 + rand.Int64N(200)},                   // Payment required
		{"api.twitter.com", "/2/tweets", "application/json", 429, 50 + rand.Int64N(100)},                             // Rate limited
		{"api.openai.com", "/v1/chat/completions", "application/json", 503, 200 + rand.Int64N(500)},                  // Service unavailable
		{"api.twilio.com", "/2010-04-01/Accounts", "application/x-www-form-urlencoded", 401, 150 + rand.Int64N(300)}, // Unauthorized
		{"api.sendgrid.com", "/v3/mail/send", "application/json", 400, 100 + rand.Int64N(200)},                       // Bad request
	}

	userAgents := []string{
		"MyApp/1.0.0 (com.mycompany.myapp; build:123)",
		"Python/3.9 requests/2.28.1",
		"Node.js/18.17.0 axios/1.4.0",
		"Go-http-client/1.1",
		"curl/7.68.0",
		"PostmanRuntime/7.32.3",
		"Java/11.0.16 OkHttp/4.10.0",
		"PHP/8.1.0 GuzzleHttp/7.5.0",
	}

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	// 85% success, 15% errors (higher error rate for external APIs)
	var selectedEndpoint struct {
		domain      string
		path        string
		contentType string
		status      int
		duration    int64
		bytes       int64
	}

	if rand.Float64() < 0.85 {
		r := apiEndpoints[rand.IntN(len(apiEndpoints))]
		selectedEndpoint = struct {
			domain      string
			path        string
			contentType string
			status      int
			duration    int64
			bytes       int64
		}{r.domain, r.path, r.contentType, r.status, r.duration, r.bytes}
	} else {
		r := errorEndpoints[rand.IntN(len(errorEndpoints))]
		selectedEndpoint = struct {
			domain      string
			path        string
			contentType string
			status      int
			duration    int64
			bytes       int64
		}{r.domain, r.path, r.contentType, r.status, r.duration, 0}
	}

	// Generate random IDs
	requestID := generateRandomID()
	connectionID := generateRandomID()
	endpointID := generateRandomID()

	// Build full URL
	fullURL := "https://" + selectedEndpoint.domain + selectedEndpoint.path

	return Request{
		ConnectionID: connectionID,
		EndpointId:   endpointID,
		RequestId:    requestID,
		Timestamp:    time.Now(),
		Direction:    "outbound",
		Url:          fullURL,
		URLPath:      selectedEndpoint.path,
		Method:       methods[rand.IntN(len(methods))],
		Status:       selectedEndpoint.status,
		Duration:     selectedEndpoint.duration,
		ContentType:  selectedEndpoint.contentType,
		Category:     "api",
		Agent:        userAgents[rand.IntN(len(userAgents))],
		WrBytes:      rand.Int64N(2048) + 256, // Larger request bodies for APIs
		RdBytes:      selectedEndpoint.bytes,
	}
}

func generateRandomID() string {
	bytes := make([]byte, 8)
	cryptorand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Convert main.Request to templates.Request
func convertToTemplateRequests(requests []Request) []templates.Request {
	result := make([]templates.Request, len(requests))
	for i, r := range requests {
		result[i] = templates.Request{
			ConnectionID:          r.ConnectionID,
			EndpointId:            r.EndpointId,
			RequestId:             r.RequestId,
			TLSProbeTypesDetected: r.TLSProbeTypesDetected,
			TLSIntrospected:       r.TLSIntrospected,
			Tags:                  r.Tags,
			Timestamp:             r.Timestamp,
			Direction:             r.Direction,
			Url:                   r.Url,
			URLPath:               r.URLPath,
			Method:                r.Method,
			Status:                r.Status,
			Duration:              r.Duration,
			ContentType:           r.ContentType,
			Category:              r.Category,
			Agent:                 r.Agent,
			WrBytes:               r.WrBytes,
			RdBytes:               r.RdBytes,
			AuthTokenMask:         r.AuthTokenMask,
			AuthTokenHash:         r.AuthTokenHash,
			AuthTokenSource:       r.AuthTokenSource,
			AuthTokenType:         r.AuthTokenType,
		}
	}
	return result
}

// Start request generator goroutine
func startRequestGenerator() {
	go func() {
		for {
			// Generate a request every 1-3 seconds
			time.Sleep(time.Duration(1000+rand.Int64N(2000)) * time.Millisecond)

			// Add request to list
			newRequest := generateRandomRequest()
			requestsMu.Lock()
			requests = append(requests, newRequest)
			// Keep only last 100 requests
			if len(requests) > 100 {
				requests = requests[len(requests)-100:]
			}
			requestsMu.Unlock()

			// Broadcast to SSE clients
			go broadcastRequestChanges()
		}
	}()
}
