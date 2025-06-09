# KochoCORS

A powerful, configurable HTTP proxy server written in Go with intelligent CORS origin detection, authentication, rate limiting, and comprehensive security features.

## Default Configuration: Intentionally Vulnerable/Permissive

> **⚠️ Warning**  
> The default configuration is designed for development and testing purposes and is intentionally permissive. It is **not recommended for production environments** due to potential security vulnerabilities.

## Quick Start

### Basic Usage (Development Mode)
```bash
# Start with default settings (no authentication, all domains allowed)
go run main.go

# Or use the compiled binary
./proxy.exe
```

### Common Usage Examples
```bash
# Development with debugging
go run main.go -debug -port=8080

# Secure mode with authentication
go run main.go -auth-key=secret123 -allowed-domains=api.example.com

# Production-like configuration
go run main.go -auth-key=prod-key -rate-limit=60 -allowed-domains=myapi.com
```

## Making Requests

```bash
# Basic proxy request (no auth required by default)
curl "http://localhost:3000/proxy?url=https://api.github.com/users/octocat"

# With authentication (when auth-key is configured)
curl "http://localhost:3000/proxy?url=https://api.example.com/data&key=secret123"

# Health check
curl "http://localhost:3000/ping"
```

## Configuration Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-port` | Port to listen on | `3000` |
| `-debug` | Enable debug logging | `false` |
| `-auth-key` | Authentication key required | `""` (no auth) |
| `-rate-limit` | Requests per minute | `0` (unlimited) |
| `-allowed-domains` | Comma-separated allowed domains | `""` (all domains) |
| `-default-origin` | Default CORS origin | `"*"` |
| `-insecure-tls` | Skip TLS certificate verification | `false` |

## Build and Run

```bash
# Install dependencies
go mod tidy

# Build executable
go build -o proxy.exe

# Run with flags
./proxy.exe -debug -port=8080
```

## Full Documentation

For comprehensive documentation including:
- Dynamic CORS origin detection details
- Security features and best practices
- Production deployment guidance
- Troubleshooting and debugging
- API reference and examples

[See Complete Documentation](./docs/README.md)

## Key Features

- Dynamic CORS Origin Detection - Intelligently detects appropriate CORS origins
- Authentication Support - Optional API key authentication
- Rate Limiting - Configurable request rate limiting
- Domain Whitelisting - Restrict proxy targets to specific domains
- Header Forwarding - Forwards ALL headers (except problematic ones)
- TLS Support - Full HTTPS support with optional cert validation bypass
- Debug Logging - Comprehensive logging for troubleshooting
- Thread-Safe - Handles concurrent requests safely

## What This Proxy Does

- Custom Header Forwarding: Forwards ALL incoming headers except 'host' and 'access-control-*'
- Concurrent Requests: Thread-safe handling of multiple simultaneous requests
- Dynamic CORS: Intelligent origin detection based on request headers
- IP Forwarding: Does NOT preserve original client IP
- Redirect Following: Does NOT automatically follow redirects

## **License**

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---
**Note**: Always secure your configuration for production environments.
