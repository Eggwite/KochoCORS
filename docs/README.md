# KochoCORS Docs

### **IMPORTANT NOTICE: Default Configuration**

> **⚠️ Warning**  
> The default configuration is designed for development and testing purposes and is intentionally permissive. It is **not recommended for production environments** due to potential security vulnerabilities.

### Default Settings:
- **No Authentication**: Accepts all requests without authentication.
- **Open Domain Access**: Proxies requests to any URL without restrictions.
- **Unlimited Rate Limiting**: No limits on request frequency.
- **Permissive CORS**: Allows requests from any origin (`*`).
- **Header Forwarding**: Forwards all incoming headers except problematic ones.
- **No IP Preservation**: Does not forward the original client IP.

## **Features**

### Supported Features:
- **Dynamic CORS Detection**: Automatically detects and sets appropriate CORS origins.
- **Authentication**: Optional authentication via the `X-KochoCORS-Auth-Token` header.
- **Rate Limiting**: Configurable request limits per minute.
- **Domain Whitelisting**: Restricts proxy access to specified domains.
- **TLS Support**: Includes HTTPS support with optional certificate validation bypass.
- **Debug Logging**: Provides detailed logs for troubleshooting.
- **Concurrent Request Handling**: Fully thread-safe with support for multiple simultaneous requests.
- **Configurable Redirect Following**: Choose whether the proxy should follow HTTP redirects from the target server or return them to the client.

### Limitations:
> **⚠️ Note**  
> - **IP Forwarding**: Does not preserve or forward the original client IP.  
> - **Redirect Handling**: By default, follows redirects. Can be configured to return redirect responses as-is.

## **Getting Started**

### Basic Usage
```bash
# Start with default settings
go run main.go

# Using the compiled binary
./proxy.exe
```

### Secure Configuration Examples
```bash
# Enable authentication and domain restrictions
go run main.go -auth-key=secret123 -allowed-domains=api.example.com,example.com

# Development mode with debugging enabled
go run main.go -debug -port=8080

# Production-like mode with rate limiting
go run main.go -auth-key=prod-secret -rate-limit=60 -allowed-domains=myapi.com
```

## **Configuration Options**

### Command-Line Flags
| Flag              | Description                        | Default Value | Example                          |
|-------------------|------------------------------------|---------------|----------------------------------|
| `-port`           | Port to listen on                 | `3000`        | `-port=8080`                    |
| `-debug`          | Enable debug logging              | `false`       | `-debug`                        |
| `-auth-key`       | Authentication key                | `""`          | `-auth-key=secret123`           |
| `-rate-limit`     | Requests per minute               | `0`           | `-rate-limit=60`                |
| `-allowed-domains`| Comma-separated allowed domains   | `""`          | `-allowed-domains=api.com,app.com` |
| `-default-origin` | Default CORS origin               | `"*"`         | `-default-origin=https://myapp.com` |
| `-insecure-tls`   | Skip TLS certificate validation   | `false`       | `-insecure-tls`                 |
| `-follow-redirects` | Follow HTTP redirects from target | `true`        | `-follow-redirects=false`       |

### Environment Variables
```bash
# Example environment variable configuration
export PORT=8080
export DEBUG=true
export AUTH_KEY=secret123
export RATE_LIMIT=60
export ALLOWED_DOMAINS=api.example.com,example.com
export DEFAULT_ORIGIN=https://myapp.com
export INSECURE_TLS=true
export FOLLOW_REDIRECTS=false
```

### .env File Support
Create a `.env` file in the project directory:
```env
PORT=8080
DEBUG=true
AUTH_KEY=secret123
RATE_LIMIT=60
ALLOWED_DOMAINS=api.example.com,example.com
DEFAULT_ORIGIN=https://myapp.com
INSECURE_TLS=false
FOLLOW_REDIRECTS=true
```

**Configuration Precedence**:
1. Command-line flags (highest priority)
2. Environment variables
3. `.env` file
4. Default values (lowest priority)

## **Request Format**

### Example Requests

#### Basic Request
```bash
curl "http://localhost:3000/proxy?url=https://api.github.com/users/octocat"
```

#### Authenticated Request
```bash
curl -H "X-KochoCORS-Auth-Token: secret123" "http://localhost:3000/proxy?url=https://api.example.com/data"
```

#### POST Request with Payload
```bash
curl -X POST \
    -H "X-KochoCORS-Auth-Token: secret123" \
    "http://localhost:3000/proxy?url=https://api.example.com/users" \
    -d '{"name": "John Doe"}'
```

#### JavaScript Example
```javascript
fetch('http://localhost:3000/proxy?url=https://api.example.com/data', {
    method: 'GET',
    headers: {
        'X-KochoCORS-Auth-Token': 'secret123'
    }
})
.then(response => response.json())
.then(data => console.log(data));
```

## **Security**

### Authentication
> **🔒 Security**  
> - Optional authentication via `key` query parameter.  
> - Returns `401 Unauthorized` for invalid or missing keys.

### Domain Whitelisting
Restrict proxy access to specific domains:
```bash
go run main.go -allowed-domains=api.example.com,secure-api.com
```

### Rate Limiting
Limit the number of requests per minute:
```bash
go run main.go -rate-limit=60
```

### TLS Certificate Validation
Enable or bypass TLS certificate validation:
```bash
# Production mode (validates certificates)
go run main.go

# Development mode (skips validation)
go run main.go -insecure-tls
```

### Header Security
> **🔒 Note**  
> - Strips `host` and `access-control-*` headers to prevent conflicts.  
> - Forwards all other headers for compatibility.

## **Development**

### Building the Project
```bash
# Build for the current platform
go build -o proxy.exe

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o proxy-linux
GOOS=darwin GOARCH=amd64 go build -o proxy-macos
```

### Dependencies
```bash
go mod init proxy-server
go get github.com/joho/godotenv
go get golang.org/x/time/rate
```

### Project Structure
```
learning-go/
├── main.go              # Main proxy server code
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── .env.example         # Example environment configuration
├── README.md            # Basic project readme
├── proxy.exe            # Compiled binary (Windows)
└── docs/
        └── README.md        # Documentation
```

## **Production Deployment**

### Recommended Settings
```bash
go run main.go \
    -auth-key=STRONG_SECRET_KEY_HERE \
    -allowed-domains=your-api.com,trusted-partner.com \
    -rate-limit=100 \
    -port=8080 \
    -default-origin=https://your-frontend.com
```

### Security Checklist
> **🔒 Security Checklist**  
> - Use a strong authentication key (`-auth-key`).  
> - Configure domain whitelisting (`-allowed-domains`).  
> - Set appropriate rate limits (`-rate-limit`).  
> - Specify a secure CORS origin (`-default-origin`).  
> - Enable TLS certificate validation (remove `-insecure-tls`).  
> - Deploy behind a reverse proxy (e.g., nginx, caddy).  

## **Troubleshooting**

### Common Issues

#### CORS Errors
> **⚠️ Note**  
> - Ensure `Origin` or `Referer` headers are sent.  
> - Verify dynamic origin detection is working.  
> - Consider setting an explicit `-default-origin`.

#### Authentication Failures
> **⚠️ Note**  
> - Verify the `auth-key` matches the server configuration.  
> - Include the `auth-key` in the request query.

#### Domain Blocked
> **⚠️ Note**  
> - Add the domain to the `-allowed-domains` flag.  
> - Ensure the domain name matches exactly.

#### Rate Limiting
> **⚠️ Note**  
> - Reduce request frequency or increase the `-rate-limit` value.

---
**Note**: Always secure your configuration for production environments.
