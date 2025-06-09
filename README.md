# KochoCORS: Your Smart & Secure CORS Proxy

**Unlock seamless cross-origin requests with KochoCORS!** A blazing-fast, highly configurable CORS proxy server built in Go. Features intelligent CORS origin detection, robust authentication, rate limiting, domain whitelisting, and more, putting you in control of your API interactions.

> **⚠️ Critical Security Warning: Default Configuration is Permissive!**  
> The default settings are for **development and testing only** and are intentionally insecure (no auth, all domains allowed). **Always customize environment variables or flags for production to secure your proxy.**

## Instantly Deploy KochoCORS

Get started in seconds! Deploy KochoCORS to your favorite platform with a single click:

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/Eggwite/KochoCORS/)
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/Eggwite/KochoCORS/)
[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/Eggwite/KochoCORS/)

## Quick Start & Usage

### Run with Docker (from Docker Hub)

First, ensure you have Docker installed. You can pull the pre-built image from Docker Hub:

```bash
docker pull eggwite/kocho-cors:latest
```

Then, run the Docker container:
```bash
# Example with port mapping and environment variables for configuration
docker run -d -p 8080:3000 \\
  -e PORT="3000" \\
  -e AUTH_KEY="your_secure_auth_key" \\
  -e ALLOWED_DOMAINS="api.example.com,another.api.com" \\
  -e RATE_LIMIT="60" \\
  -e DEFAULT_ORIGIN="*" \\
  -e INSECURE_TLS="false" \\
  -e FOLLOW_REDIRECTS="true" \\
  -e DEBUG="false" \\
  eggwite/kocho-cors:latest
```
*(The application inside Docker defaults to port 3000 if `PORT` env var is not set. Adjust environment variables as needed.)*

Alternatively, if you've built the image locally (e.g., `docker build -t kocho-cors .`):
```bash
docker run -d -p 8080:3000 \\
  -e PORT="3000" \\
  -e AUTH_KEY="your_secure_auth_key" \\
  # ... other env vars
  kocho-cors
```
### Run Locally (Development)
Install the repository and navigate to the directory.
```bash
# Start with default, permissive settings
go run main.go

# Or use a pre-compiled binary (e.g., after `go build -o proxy.exe`)
./proxy.exe
```

### Common Usage Examples with Flags
```bash
# Development with debugging on a different port
go run main.go -debug -port=8080

# Secure mode with authentication and domain whitelisting
go run main.go -auth-key=secret123 -allowed-domains=api.example.com

# Production-like configuration with rate limiting
go run main.go -auth-key=prod-key -rate-limit=60 -allowed-domains=myapi.com
```

## Key Features at a Glance

-   **🚀 Dynamic CORS Origin Detection**: Intelligently sets `Access-Control-Allow-Origin`.
-   **🔑 Authentication Support**: Secure your proxy with an optional API key passed via the `X-KochoCORS-Auth-Token` header.
-   **⏱️ Rate Limiting**: Protect your target APIs with configurable request limits.
-   **🛡️ Domain Whitelisting**: Restrict proxy access to only approved target domains.
-   **📋 Comprehensive Header Forwarding**: Forwards most client headers to the target.
-   **🔒 TLS Support**: Full HTTPS for secure connections, with optional insecure skip for development.
-   **🐛 Debug Logging**: Detailed logs for easy troubleshooting.
-   **⚡ Thread-Safe & Concurrent**: Built to handle many requests efficiently.

## Making Requests

Once KochoCORS is running (e.g., on `http://localhost:3000` by default):

```bash
# Basic proxy request (if no auth-key is configured)
curl "http://localhost:3000/proxy?url=https://example.com"

# With authentication (when auth-key is configured)
curl -H "X-KochoCORS-Auth-Token: YOUR_AUTH_KEY" "http://localhost:3000/proxy?url=https://api.example.com/data"

# Health check endpoint
curl "http://localhost:3000/ping"
```

## What This Proxy Does (And Doesn't)

-   **Custom Header Forwarding**: Forwards ALL incoming headers from the client to the target URL, except for `Host` and `Access-Control-*` headers which are managed by the proxy.
-   **Concurrent Requests**: Designed to be thread-safe, efficiently handling multiple simultaneous requests.
-   **Dynamic CORS Handling**: Intelligently determines the `Access-Control-Allow-Origin` header based on the incoming request's `Origin` or `Referer` if `DefaultOrigin` is `*` and an `AuthKey` is set. Otherwise, uses `DefaultOrigin`.
-   **IP Forwarding**: Does **NOT** preserve or forward the original client's IP address to the target server (e.g., no `X-Forwarded-For`).
-   **Redirect Following**: Does **NOT** automatically follow redirects from the target server. It returns the redirect response (e.g., 301, 302) as-is to the client.

## Configuration Flags

All settings can be configured via command-line flags or environment variables (environment variables take precedence if set, e.g., `PORT` for `-port`, `AUTH_KEY` for `-auth-key`).

| Flag                | Environment Variable | Description                               | Default          |
| ------------------- | -------------------- | ----------------------------------------- | ---------------- |
| `-port`             | `PORT`               | Port to listen on                         | `3000`           |
| `-debug`            | `DEBUG`              | Enable debug logging (true/false)         | `false`          |
| `-auth-key`         | `AUTH_KEY`           | Authentication key required for requests (passed in `X-KochoCORS-Auth-Token` header)  | `""` (no auth)   |
| `-rate-limit`       | `RATE_LIMIT`         | Requests per minute (0 to disable)        | `0` (unlimited)  |
| `-allowed-domains`  | `ALLOWED_DOMAINS`    | Comma-separated allowed target domains    | `""` (all allowed) |
| `-default-origin`   | `DEFAULT_ORIGIN`     | Default `Access-Control-Allow-Origin`     | `*`              |
| `-insecure-tls`     | `INSECURE_TLS`       | Skip TLS certificate verification (true/false) | `false`          |
| `-follow-redirects` | `FOLLOW_REDIRECTS`   | Follow HTTP redirects from target URL (true/false) | `true`           |

## Build and Run (Manual)

```bash
# Ensure Go is installed, then install dependencies
go mod tidy

# Build the executable (e.g., proxy.exe on Windows, proxy on Linux/macOS)
go build -o proxy.exe . # Or: go build -o proxy .

# Run with desired flags
./proxy.exe -debug -port=8080
```

## Platform-Specific Deployment Notes & Docker Hub

When deploying, refer to the platform-specific configuration files in the root of this repository (`render.yaml`, `Dockerfile`, `vercel.json`, `Procfile`).

#### Render (`render.yaml`)
*   **Plan**: Defaults to `free`. Modify `plan: free` as needed.
*   **Health Check**: `healthCheckPath: /ping` is used by Render to monitor service health.
*   **Port**: Render injects a `PORT` environment variable. The `value: "10000"` in `envVars` is a suggestion; the application will use the `PORT` provided by Render.
*   **Environment Variables**: Configure `AUTH_KEY`, `ALLOWED_DOMAINS`, etc., in the Render dashboard or `render.yaml`'s `envVars`.

#### Heroku (`Procfile`)
*   The Go buildpack compiles `main.go`. `web: proxy` (or `./proxy`) runs the executable.
*   Heroku sets the `PORT` environment variable. Configure others via the dashboard/CLI.

#### Vercel (`vercel.json`)
*   Uses `@vercel/go` builder. Vercel handles routing and sets the `PORT`.
*   Configure other environment variables in the Vercel project settings.

#### Docker (`Dockerfile` details)
*   Provides a multi-stage build for a small, optimized Linux image.
*   Listens on the `PORT` env var (defaults to `3000`). `EXPOSE 3000` documents this.
*   Pass environment variables during `docker run` as shown in the Quick Start.

## Full Documentation

For comprehensive details on all features, advanced configurations, security best practices, and troubleshooting:

➡️ **[Read the Complete KochoCORS Documentation](./docs/README.md)**

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---
**Reminder**: Always secure your configuration for production environments by setting appropriate flags or environment variables.