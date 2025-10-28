package templates

import (
	"time"
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

// Old Todo struct kept for backwards compatibility
type Todo struct {
	ID   int
	Text string
	Done bool
}
