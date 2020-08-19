package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"

	logger "github.com/TykTechnologies/tyk/log"
)

var log = logger.Get()

// Run on startup.  Bootstrapping the service here
func init() {}

// Required
func main() {}

// Main method to be run by Tyk
func CertHeaderInject(w http.ResponseWriter, r *http.Request) {
	log.Info("-----------")
	log.Info("Attempting to pull Peer Cert Info")
	if len(r.TLS.PeerCertificates) > 0 {
		// Get Issuer
		r.Header.Set("X-Client-Issuer", r.TLS.PeerCertificates[0].Issuer.CommonName)
		// Get The cert Fingerprint
		r.Header.Set("X-Client-Fingerprint", getFingerprintString(sha1.Sum(r.TLS.PeerCertificates[0].Raw)))
	}
}

func getFingerprintString(fingerprint [20]byte) string {
	var buf bytes.Buffer
	for i, f := range fingerprint {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02X", f)
	}
	fmt.Printf("Fingerprint: %s", buf.String())
	return buf.String()
}
