package main

import (
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
	clientIP := getIP(r.RemoteAddr)

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

func getIP(remoteAddr string) string {
	if ip := strings.Split(remoteAddr, ":"); len(ip) > 0 {
		return ip[0]
	}
	return ""
}
