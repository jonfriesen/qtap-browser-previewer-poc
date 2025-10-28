# Dear LLM - Network Monitor Application Notes

## Quick Start Guide

This is a Go web application that simulates a Chrome DevTools Network tab interface with real-time request monitoring.

### Running the Application

Before starting any work, ensure that you install some dependencies:

- Tailwind: `npm install tailwindcss @tailwindcss/cli`
- Templ: `go install github.com/a-h/templ/cmd/templ@latest`

**Preferred Method:**

```bash
# Generate templates first (required after any .templ file changes)
templ generate .

# Build and run
go build -o server .
PORT=8080 ./server
```

**Alternative Method:**

```bash
# Direct run (slower but works)
go run main.go handlers.go sse.go
```

### Common Issues & Solutions

#### 1. Template Generation

- **Always run `templ generate .` after modifying `.templ` files**
- The application uses Go templ for HTML templates
- Templates compile to `*_templ.go` files that need regeneration

#### 2. Port Binding Issues

- Default port is 8080, but it may conflict in some environments
- Use `PORT=9090` or another port if 8080 fails
- Server binds to `0.0.0.0:PORT` so it should be accessible

#### 3. Background Process Management

- Processes may exit silently if there are template/compilation errors
- Test with `timeout 5s ./server` first to see if it starts properly
- Use `ps aux | grep server` to verify the process is running

#### 4. Static Files

- CSS files are in `static/css/`
- Server serves static files at `/static/` route
- Font size classes: `text-2xs` (10px), `text-xs` (11px), `text-sm` (12px)

### Architecture Notes

- **Backend**: Go with standard HTTP library
- **Frontend**: HTML with HTMX for dynamic updates
- **Real-time**: Server-Sent Events (SSE) for live request streaming
- **Styling**: Custom CSS with Tailwind-like utility classes
- **Templates**: Go templ for server-side rendering

### File Structure

```
├── main.go           # HTTP server setup
├── handlers.go       # Request handlers & data generation
├── sse.go           # Server-Sent Events implementation
├── templates/
│   ├── base.templ   # Base HTML template
│   ├── todo.templ   # Main network interface template
│   └── *_templ.go   # Generated Go files (auto-generated)
└── static/css/
    └── style.css    # Chrome DevTools-inspired styling
```

### Development Workflow

1. Make changes to `.templ` files or Go code
2. Run `templ generate .` if templates were modified
3. Test with `go build -o server . && ./server`
4. Access at `http://localhost:8080`
5. Use browser dev tools to inspect real-time SSE updates

### Features

- **Live Request Generation**: Creates realistic fake network requests every 1-3 seconds
- **Chrome DevTools UI**: Faithful reproduction of Network tab interface
- **Real-time Updates**: Uses SSE to stream new requests to connected clients
- **Interactive Controls**: Clear button, filter options, etc.
- **Realistic Data**: HTTP methods, status codes, file types, timing, sizes

### Debugging Tips

- Check `go mod tidy` if there are dependency issues
- Verify templates compile: `templ generate . && go build .`
- Test basic connectivity: `curl -I http://localhost:8080`
- Monitor process: `ps aux | grep server`
- Check for port conflicts: `netstat -tlnp | grep 8080`

---

_Last updated: During font size reduction implementation_
