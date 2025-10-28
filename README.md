# Network Monitor Demo

A Chrome DevTools Network tab-inspired web application built with Go using modern web technologies.

## Tech Stack

- **Go** - Backend server and logic
- **[templ](https://github.com/a-h/templ)** - Go HTML templating system
- **[Tailwind CSS](https://tailwindcss.com/)** - Dark theme styling
- **[HTMX](https://htmx.org/)** - Dynamic HTML updates
- **[HTMX SSE Extension](https://htmx.org/extensions/server-sent-events/)** - Real-time request streaming
- **[Alpine.js](https://alpinejs.dev/)** - Client-side JavaScript (ready for future enhancements)

## Quick Start

1. **Install dependencies:**

   ```bash
   go mod download
   go install github.com/a-h/templ/cmd/templ@latest
   ```

2. **Install Tailwind CSS:**

   Preferred:

   ```bash
   npm install tailwindcss @tailwindcss/cli
   ```

   ```bash
   curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64
   chmod +x tailwindcss-linux-arm64
   sudo mv tailwindcss-linux-arm64 /usr/local/bin/tailwindcss
   ```

3. **Generate templates and CSS:**

   ```bash
   templ generate
   tailwindcss -i static/css/input.css -o static/css/style.css --minify
   ```

4. **Run the application:**

   ```bash
   go run .
   ```

5. **View the network monitor:**
   Navigate to `http://localhost:8080` to see the Chrome DevTools-inspired network monitor in action!

## Development

- Templates are in `templates/` - edit `.templ` files and run `templ generate`
- Styles are in `static/css/input.css` - rebuild with tailwindcss command
- Server logic is in `main.go`, `handlers.go`, and `sse.go`
- Request generation logic is in `handlers.go`

## Project Structure

```
.
├── main.go              # Server setup and routing
├── handlers.go          # HTTP handlers and request generation
├── sse.go               # Server-Sent Events implementation
├── templates/           # Templ template files
│   ├── base.templ      # Base HTML layout with dark theme
│   └── todo.templ      # Network monitor components and request table
├── static/             # Static assets
│   └── css/
│       ├── input.css   # Tailwind input file
│       └── style.css   # Generated CSS with dark theme
└── tailwind.config.js  # Tailwind configuration
```
