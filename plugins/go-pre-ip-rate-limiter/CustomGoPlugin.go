package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/user"
)

var logger = log.Get()

func main() {}

// This will be run during Gateway startup
func init() {
	logger.Info("--- Go custom plugin init success! ---- ")
}

// IPRateLimit
func IPRateLimit(rw http.ResponseWriter, r *http.Request) {
	// Get the client IP address
	clientIP := getIP(r)

	logger.Info("Request came in from " + clientIP)

	// Create a Redis key for the IP address
	session := &user.SessionState{
		Alias: clientIP,
		Rate:  2,
		Per:   10,
		MetaData: map[string]interface{}{
			"token": clientIP,
		},
		KeyID: clientIP,
	}

	ctx.SetSession(r, session, true)
}

func getIP(r *http.Request) string {

	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, addr := range strings.Split(r.Header.Get(h), ",") {
			addr = strings.TrimSpace(addr)
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(addr)
			if !realIP.IsGlobalUnicast() {
				// bad address - we should alsp eventually check for private IP ranges
				continue
			}
			return addr
		}
	}
	// If no valid IP found, return the remote address
	return r.RemoteAddr
}
