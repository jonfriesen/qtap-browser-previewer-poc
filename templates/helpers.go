package templates

import (
	"fmt"
	"strings"
)

// Helper functions for template rendering

func getFileName(path string) string {
	if path == "" || path == "/" {
		return "(index)"
	}
	// Get last part of path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && parts[len(parts)-1] != "" {
		return parts[len(parts)-1]
	}
	return path
}

func getRequestType(contentType string) string {
	switch {
	case strings.Contains(contentType, "javascript") || strings.Contains(contentType, "application/javascript"):
		return "script"
	case strings.Contains(contentType, "css"):
		return "stylesheet"
	case strings.Contains(contentType, "html"):
		return "document"
	case strings.Contains(contentType, "image"):
		return "image"
	case strings.Contains(contentType, "json"):
		return "xhr"
	case strings.Contains(contentType, "font"):
		return "font"
	default:
		return "other"
	}
}

func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "text-green-400"
	case status >= 300 && status < 400:
		return "text-yellow-400"
	case status >= 400:
		return "text-red-400"
	default:
		return "text-gray-400"
	}
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f kB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
}

func formatTotalBytes(requests []Request, includeWritten bool) string {
	total := int64(0)
	for _, r := range requests {
		total += r.RdBytes
		if includeWritten {
			total += r.WrBytes
		}
	}
	return formatBytes(total)
}

func getMaxDuration(requests []Request) int64 {
	max := int64(0)
	for _, r := range requests {
		if r.Duration > max {
			max = r.Duration
		}
	}
	return max
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getDomain(url string) string {
	// Handle URLs with protocol
	if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	}

	// Extract domain part (everything before first slash or colon)
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}
	if idx := strings.Index(url, ":"); idx != -1 {
		url = url[:idx]
	}

	if url == "" {
		return "localhost"
	}

	return url
}
