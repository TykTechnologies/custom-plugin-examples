package main

import (
	"fmt"
	"net/http"
	"strings"
)

// setCORSHeaders writes permissive CORS headers on the response.
// These headers enable the portal's "Try It" (Stoplight Elements) to call mock endpoints
// cross-origin from the browser.
func setCORSHeaders(rw http.ResponseWriter, cfg corsConfig) {
	rw.Header().Set("Access-Control-Allow-Origin", strings.Join(cfg.AllowOrigins, ", "))
	rw.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
	rw.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
	rw.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
}

// handleCORSPreflight handles OPTIONS requests for mock endpoints.
// Returns true if the request was handled (caller should return immediately).
func handleCORSPreflight(rw http.ResponseWriter, r *http.Request, cfg mockConfig) bool {
	if r.Method != http.MethodOptions {
		return false
	}

	setCORSHeaders(rw, cfg.CORS)
	rw.WriteHeader(http.StatusNoContent)
	return true
}
