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
		// Requested Hostname, same as nginx "$host"
		r.Header.Set("x-client-requested-servername", r.TLS.ServerName)
		log.Info("x-client-requested-servername: " + r.TLS.ServerName)
		// Client Cert Issuer
		r.Header.Set("X-Client-Issuer", r.TLS.PeerCertificates[0].Issuer.CommonName)
		log.Info("X-Client-Issuer: " + r.TLS.PeerCertificates[0].Issuer.CommonName)
		// Client Cert Fingerprint
		r.Header.Set("X-Client-Fingerprint", getFingerprintString(sha1.Sum(r.TLS.PeerCertificates[0].Raw)))
		log.Info("X-Client-Fingerprint: " + getFingerprintString(sha1.Sum(r.TLS.PeerCertificates[0].Raw)))
	} else {
		log.Info("No Peer certificates found.")
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
