package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

type Config struct {
	AllowedDomains  []string
	RateLimit       int
	Port            string
	AuthKey         string
	Debug           bool
	InsecureTLS     bool
	DefaultOrigin   string
	FollowRedirects bool
}

var (
	flagPort            = flag.String("port", "", "port to listen on")
	flagDebug           = flag.Bool("debug", false, "enable debug logging")
	flagInsecureTLS     = flag.Bool("insecure-tls", false, "skip TLS certificate verification")
	flagDefaultOrigin   = flag.String("default-origin", "*", "default origin for CORS")
	flagAuthKey         = flag.String("auth-key", "", "authentication key required for proxy requests, passed via X-KochoCORS-Auth-Token header")
	flagRateLimit       = flag.Int("rate-limit", 0, "rate limit per minute (0 to disable)")
	flagAllowedDomains  = flag.String("allowed-domains", "", "comma-separated list of allowed domains")
	flagFollowRedirects = flag.Bool("follow-redirects", true, "follow HTTP redirects from the target URL")
)

var appConfig Config
var limiter *rate.Limiter
var once sync.Once

// loadConfig initializes the application configuration from flags, environment variables, and defaults.
func loadConfig() {
	once.Do(func() {
		flag.Parse()

		err := godotenv.Load()
		if err != nil {
			logger("No .env file found, using default or environment variable values")
		}

		domains := os.Getenv("ALLOWED_DOMAINS")
		if *flagAllowedDomains != "" {
			domains = *flagAllowedDomains
		}
		if domains != "" {
			appConfig.AllowedDomains = strings.Split(domains, ",")
		}

		rateLimitStr := os.Getenv("RATE_LIMIT")
		if *flagRateLimit > 0 {
			appConfig.RateLimit = *flagRateLimit
		} else if rateLimitStr != "" {
			val, err := strconv.Atoi(rateLimitStr)
			if err == nil {
				appConfig.RateLimit = val
			}
		}

		appConfig.Port = os.Getenv("PORT")
		if *flagPort != "" {
			appConfig.Port = *flagPort
		}
		if appConfig.Port == "" {
			appConfig.Port = "3000"
		}

		appConfig.AuthKey = os.Getenv("AUTH_KEY")
		if *flagAuthKey != "" {
			appConfig.AuthKey = *flagAuthKey
		}

		appConfig.Debug = *flagDebug
		if os.Getenv("DEBUG") == "true" {
			appConfig.Debug = true
		}

		appConfig.InsecureTLS = *flagInsecureTLS
		if os.Getenv("INSECURE_TLS") == "true" {
			appConfig.InsecureTLS = true
		}

		appConfig.DefaultOrigin = *flagDefaultOrigin
		if origin := os.Getenv("DEFAULT_ORIGIN"); origin != "" {
			appConfig.DefaultOrigin = origin
		}

		appConfig.FollowRedirects = *flagFollowRedirects
		if os.Getenv("FOLLOW_REDIRECTS") == "false" {
			appConfig.FollowRedirects = false
		} else if os.Getenv("FOLLOW_REDIRECTS") == "true" {
			appConfig.FollowRedirects = true
		}

		if appConfig.RateLimit > 0 {
			limiter = rate.NewLimiter(rate.Limit(appConfig.RateLimit)/60.0, 1)
		}

		logger(fmt.Sprintf("Configuration loaded: %+v", appConfig))
	})
}

func logger(v ...interface{}) {
	if appConfig.Debug {
		log.Println(v...)
	}
}

func validateConfig() {
	if appConfig.Port == "" {
		log.Fatal("Port cannot be empty")
	}
}

// getCORSOrigin determines the appropriate CORS origin based on the request and configuration.
func getCORSOrigin(r *http.Request) string {
	if appConfig.DefaultOrigin != "*" {
		return appConfig.DefaultOrigin
	}

	if appConfig.AuthKey != "" {
		if origin := r.Header.Get("Origin"); origin != "" {
			return origin
		}

		if referer := r.Header.Get("Referer"); referer != "" {
			if parsedURL, err := url.Parse(referer); err == nil {
				return parsedURL.Scheme + "://" + parsedURL.Host
			}
		}

		return appConfig.DefaultOrigin
	}

	return "*"
}

// Define routes
func main() {
	loadConfig()
	validateConfig()

	// Ha ha ping pong
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "pong",
		})
	})

	// Proxy endpoint
	http.HandleFunc("/proxy", handleProxyRequest)

	log.Printf("Server running on http://localhost:%s", appConfig.Port)
	log.Printf("Proxy endpoint: http://localhost:%s/proxy?url=TARGET_URL", appConfig.Port)
	if appConfig.AuthKey != "" {
		log.Printf("Authentication required: pass token in X-KochoCORS-Auth-Token header")
	}

	if err := http.ListenAndServe(":"+appConfig.Port, nil); err != nil {
		log.Fatal(err)
	}
}

// handleProxyRequest is main logical loop
func handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	if appConfig.AuthKey != "" {
		providedKey := r.Header.Get("X-KochoCORS-Auth-Token")
		if providedKey != appConfig.AuthKey {
			http.Error(w, "Unauthorized - Invalid or missing X-KochoCORS-Auth-Token header", http.StatusUnauthorized)
			return
		}
	}

	if limiter != nil {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
	}

	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "Missing url query parameter", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(targetURL)
	if err != nil {
		http.Error(w, "Invalid URL provided", http.StatusBadRequest)
		return
	}

	// Check if the domain is allowed
	if len(appConfig.AllowedDomains) > 0 {
		isAllowed := false
		for _, domain := range appConfig.AllowedDomains {
			if domain != "" && strings.HasSuffix(parsedURL.Hostname(), strings.TrimSpace(domain)) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			http.Error(w, "Domain not allowed: "+parsedURL.Hostname(), http.StatusForbidden)
			return
		}
	}

	// Set CORS headers
	corsOrigin := getCORSOrigin(r)
	w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, TRACE, COPY, LINK")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if corsOrigin != "*" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request to target URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	for name, values := range r.Header {
		lowerName := strings.ToLower(name)
		if lowerName == "host" || strings.HasPrefix(lowerName, "access-control-") {
			continue
		}
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	client := &http.Client{}

	if !appConfig.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			// If we don't want to follow redirects, return the last response,
			// which includes the redirect status code and Location header.
			return http.ErrUseLastResponse
		}
	}

	if parsedURL.Scheme == "https" || appConfig.InsecureTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: appConfig.InsecureTLS,
			},
		}
	}

	resp, err := client.Do(req) // Network request to target URL
	if err != nil {
		http.Error(w, "Failed to fetch target URL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Redirect response handling is now managed by the client's CheckRedirect policy
	// if resp.StatusCode >= 300 && resp.StatusCode < 400 {
	// 	http.Error(w, "Redirect detected: "+resp.Status, http.StatusBadRequest)
	// 	return
	// }

	// Copy response headers and body to the client
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
