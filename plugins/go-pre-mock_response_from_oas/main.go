// MockResponder — Go Native Pre-Plugin for Tyk Gateway
//
// Returns OAS example responses for mock requests, enabling portal "Try It"
// without credentials or a live backend. Runs as pre-plugin BEFORE auth.
//
// Trigger: X-Mock: true header or Host: mock.* subdomain.
// When triggered, short-circuits the entire middleware chain (auth, proxy, response plugins).
// When NOT triggered, returns immediately without writing — lets normal auth chain proceed.
//
// Example priority matches the gateway's built-in mockResponse middleware:
//   1. Direct inline example (media.Example)
//   2. Named example by X-Tyk-Accept-Example-Name header (sorted first if not specified)
//   3. Schema-based fallback via oas.ExampleExtractor
//
// Configuration via pluginConfig.data.value.mock_responder in API definition.
package main

import (
	"context"
	"net/http"
	"strings"

	tykctx "github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
)

var logger = log.Get()

func init() {
	logger.Info("[MockResponder] Plugin loaded")
}

// MockResponder is the exported pre-plugin entry point.
// It intercepts mock requests and returns OAS example responses.
// Non-mock requests pass through without any modification.
//
// Special handling for CORS: OPTIONS preflight requests are handled even without
// X-Mock header, because browsers send OPTIONS automatically without custom headers.
// This is essential for portal "Try It" (Stoplight Elements) to work cross-origin.
func MockResponder(rw http.ResponseWriter, r *http.Request) {
	// CORS preflight: browsers send OPTIONS WITHOUT custom headers (no X-Mock).
	// Must handle before the isMockRequest check to support portal "Try It".
	if r.Method == http.MethodOptions && isCORSPreflight(r) {
		raw := readPluginConfig(r)
		cfg := parseMockConfig(raw)
		if cfg.Enabled {
			handleCORSPreflight(rw, r, cfg)
			suppressAnalytics(r)
			return
		}
	}

	// Fast path: not a mock request — return immediately, let auth chain proceed
	if !isMockRequest(r) {
		return
	}

	// Get OAS definition from request context (injected by gateway before pre-plugins)
	oasDef := tykctx.GetOASDefinition(r)
	if oasDef == nil {
		logger.Debug("[MockResponder] OAS definition not in context, passing through")
		return
	}

	// Read per-API config from pluginConfig.data.value.mock_responder
	raw := readPluginConfig(r)
	cfg := parseMockConfig(raw)
	if !cfg.Enabled {
		return
	}

	// Handle mock OPTIONS with X-Mock header (explicit mock CORS)
	if handleCORSPreflight(rw, r, cfg) {
		suppressAnalytics(r)
		return
	}

	// Find matching operation and respond with example
	if !findAndRespond(rw, r, oasDef, cfg) {
		return
	}

	suppressAnalytics(r)
}

// isCORSPreflight checks if the request is a CORS preflight from a browser.
// Browsers send Access-Control-Request-Method on preflight OPTIONS.
func isCORSPreflight(r *http.Request) bool {
	return r.Header.Get("Access-Control-Request-Method") != ""
}

// isMockRequest checks if the request should be handled as a mock.
// Checks X-Mock header (case-insensitive) and Host subdomain.
func isMockRequest(r *http.Request) bool {
	if strings.EqualFold(r.Header.Get("X-Mock"), "true") {
		return true
	}
	host := r.Host
	if idx := strings.IndexByte(host, ':'); idx != -1 {
		host = host[:idx]
	}
	return strings.HasPrefix(host, "mock.") || strings.HasPrefix(host, "mock-")
}

// suppressAnalytics sets DoNotTrackThisEndpoint in the request context.
// This prevents RecordHit() from writing mock traffic to analytics.
// Uses the same context mutation pattern as the gateway's internal middleware
// (httpctx.Set does *r = *r.WithContext(newCtx) for in-place mutation).
func suppressAnalytics(r *http.Request) {
	newCtx := context.WithValue(r.Context(), tykctx.DoNotTrackThisEndpoint, true)
	*r = *r.WithContext(newCtx)
}

func main() {}
