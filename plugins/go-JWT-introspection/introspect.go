package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/TykTechnologies/tyk/log"
)

var (
	introspectionClient      *http.Client
	introspectLogger         = log.Get()
	introspectionEndpoint    = ""
	authorizationHeaderName  = "Authorization"
	authorizationHeaderValue = ""
	authMode                 = "basic" // "basic" (default) or "body"
	clientIDEnv              = ""
	clientSecretEnv          = ""
)

// IntrospectResponse represents the minimal fields used from an RFC 7662
// token introspection response. Keycloak returns these fields among others.
type IntrospectResponse struct {
	Active bool    `json:"active"`
	Scope  *string `json:"scope,omitempty"`
}

func init() {
	introspectLogger.Info("Plugin initializing")
	introspectionClient = &http.Client{}

	introspectionEndpoint = os.Getenv("OAUTH2_INTROSPECT_ENDPOINT")
	authorizationHeaderValue = os.Getenv("OAUTH2_INTROSPECT_AUTHORIZATION")
	clientIDEnv = os.Getenv("OAUTH2_INTROSPECT_CLIENT_ID")
	clientSecretEnv = os.Getenv("OAUTH2_INTROSPECT_CLIENT_SECRET")
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("OAUTH2_INTROSPECT_AUTH_MODE"))); v != "" {
		authMode = v
	}

	// If explicit Authorization header is not provided, try to build a Basic header
	// from OAUTH2_INTROSPECT_CLIENT_ID and OAUTH2_INTROSPECT_CLIENT_SECRET
	if authorizationHeaderValue == "" && authMode == "basic" {
		if clientIDEnv != "" && clientSecretEnv != "" {
			creds := clientIDEnv + ":" + clientSecretEnv
			b64 := base64.StdEncoding.EncodeToString([]byte(creds))
			authorizationHeaderValue = "Basic " + b64
			introspectLogger.Info("oauth2-introspect: using Basic auth built from client credentials env vars (client_id set, secret set)")
		} else {
			if clientIDEnv == "" {
				introspectLogger.Debug("oauth2-introspect: OAUTH2_INTROSPECT_CLIENT_ID not set")
			}
			if clientSecretEnv == "" {
				introspectLogger.Debug("oauth2-introspect: OAUTH2_INTROSPECT_CLIENT_SECRET not set")
			}
		}
	} else {
		scheme := ""
		if strings.HasPrefix(strings.ToLower(authorizationHeaderValue), "basic ") {
			scheme = "Basic"
		} else if strings.HasPrefix(strings.ToLower(authorizationHeaderValue), "bearer ") {
			scheme = "Bearer"
		} else if authorizationHeaderValue != "" {
			scheme = "Custom"
		}
		introspectLogger.Info("oauth2-introspect: using explicit Authorization header (scheme=" + scheme + ") from OAUTH2_INTROSPECT_AUTHORIZATION")
	}

	if introspectionEndpoint == "" {
		introspectLogger.Warn("oauth2-introspect: OAUTH2_INTROSPECT_ENDPOINT is not set; plugin will fail to call introspection endpoint")
	} else {
		introspectLogger.Info("oauth2-introspect: introspection endpoint set to: " + introspectionEndpoint)
	}

	// Log auth mode choice (no secrets)
	introspectLogger.Info("oauth2-introspect: auth mode=" + authMode)
	introspectLogger.Info("Plugin initialized")
}

func OAuth2Introspect(w http.ResponseWriter, r *http.Request) {
	introspectLogger.Info("Starting INTROSPEKTION")

	bearerToken := accessTokenFromRequest(r)
	if bearerToken == "" {
		introspectLogger.Info("no bearer token found in request")
		writeUnauthorizedJSON(w, "Bearer token required")
		return
	}
	introspectLogger.Info("Bearer token: " + bearerToken)

	introspectLogger.Debug("oauth2-introspect: bearer token detected (length=" + strconv.Itoa(len(bearerToken)) + ")")

	data := url.Values{}
	data.Set("token", bearerToken)
	data.Set("token_type_hint", "access_token")
	introspectionReq, err := http.NewRequest(http.MethodPost, introspectionEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		introspectLogger.Errorf("unable to create new request %s", err.Error())
		writeJSONError(w, http.StatusInternalServerError, "introspection_request_error", "Failed to build introspection request")
		return
	}

	// Attach authentication for AS call based on mode
	// Recompute body if we need to add client credentials to it
	bodyEncoded := data.Encode()
	if authMode == "body" {
		if clientIDEnv == "" || clientSecretEnv == "" {
			introspectLogger.Errorf("oauth2-introspect: auth_mode=body but client_id/secret not set; set OAUTH2_INTROSPECT_CLIENT_ID and OAUTH2_INTROSPECT_CLIENT_SECRET")
			writeJSONError(w, http.StatusInternalServerError, "introspection_auth_error", "Authorization server rejected introspection credentials")
			return
		}
		data.Set("client_id", clientIDEnv)
		data.Set("client_secret", clientSecretEnv)
		bodyEncoded = data.Encode()
		introspectionReq.Body = io.NopCloser(strings.NewReader(bodyEncoded))
		introspectionReq.ContentLength = int64(len(bodyEncoded))
		introspectLogger.Info("oauth2-introspect: using auth_mode=body (client_id set, secret set)")
	} else {
		// Default: send Authorization header
		introspectLogger.Info("oauth2-introspect: using Authorization header: " + authorizationHeaderValue)
		introspectionReq.Header.Add(authorizationHeaderName, authorizationHeaderValue)
	}
	introspectionReq.Header.Add("content-length", strconv.Itoa(len(bodyEncoded)))
	introspectionReq.Header.Add("content-type", "application/x-www-form-urlencoded")

	// Log outgoing request details without sensitive info
	scheme := ""
	if authMode == "body" {
		scheme = "Body(client_id+client_secret)"
	} else if strings.HasPrefix(strings.ToLower(authorizationHeaderValue), "basic ") {
		scheme = "Basic"
	} else if strings.HasPrefix(strings.ToLower(authorizationHeaderValue), "bearer ") {
		scheme = "Bearer"
	} else if authorizationHeaderValue != "" {
		scheme = "Custom"
	}
	introspectLogger.Info("oauth2-introspect: POST " + introspectionEndpoint + " (auth scheme=" + scheme + ", body_len=" + strconv.Itoa(len(bodyEncoded)) + ")")
	introspectionRes, err := introspectionClient.Do(introspectionReq)
	if err != nil {
		introspectLogger.Info("tyk cannot connect to the authorization server %s\n", err.Error())
		writeJSONError(w, http.StatusInternalServerError, "introspection_connect_error", "Cannot reach authorization server")
		return
	}
	if introspectionRes.StatusCode == http.StatusUnauthorized {
		// If the protected resource uses OAuth 2.0 client credentials to
		//   authenticate to the introspection endpoint and its credentials are
		//   invalid, the authorization server responds with an HTTP 401
		//   (Unauthorized) as described in Section 5.2 of OAuth 2.0 [RFC6749].
		// If the protected resource uses an OAuth 2.0 bearer token to authorize
		//   its call to the introspection endpoint and the token used for
		//   authorization does not contain sufficient privileges or is otherwise
		//   invalid for this request, the authorization server responds with an
		//   HTTP 401 code as described in Section 3 of OAuth 2.0 Bearer Token
		//   Usage [RFC6750].
		// Read body to include in logs (helps diagnose invalid client credentials)
		b, _ := io.ReadAll(introspectionRes.Body)
		snippet := string(b)
		if len(snippet) > 256 {
			snippet = snippet[:256] + "..."
		}
		wwwAuth := introspectionRes.Header.Get("WWW-Authenticate")
		introspectLogger.Info("oauth2-introspect: 401 from AS. WWW-Authenticate=%s, body=%s", wwwAuth, snippet)
		introspectLogger.Info("tyk is not authorized to perform introspection")
		// Close explicitly since we consumed the body and will return
		introspectionRes.Body.Close()
		writeJSONError(w, http.StatusInternalServerError, "introspection_auth_error", "Authorization server rejected introspection credentials")
		return
	}
	defer introspectionRes.Body.Close()

	body, err := io.ReadAll(introspectionRes.Body)
	if err != nil {
		introspectLogger.Errorf("unable to read response body from authorization server %s", err.Error())
		writeJSONError(w, http.StatusInternalServerError, "introspection_read_error", "Failed to read response from authorization server")
		return
	}

	irObj := &IntrospectResponse{}
	err = json.Unmarshal(body, irObj)
	if err != nil {
		introspectLogger.Errorf("unable to read json response from authorization server %s", err.Error())
		writeJSONError(w, http.StatusInternalServerError, "introspection_parse_error", "Invalid JSON from authorization server")
		return
	}
	if irObj.Active == false {
		// If the introspection call is properly authorized but the token is not
		//   active, does not exist on this server, or the protected resource is
		//   not allowed to introspect this particular token, then the
		//   authorization server MUST return an introspection response with the
		//   "active" field set to "false".  Note that to avoid disclosing too
		//   much of the authorization server's state to a third party, the
		//   authorization server SHOULD NOT include any additional information
		//   about an inactive token, including why the token is inactive.
		introspectLogger.Debug("oauth2-introspect: access_token is inactive (introspection returned active=false)")
		writeUnauthorizedJSON(w, "Access token is inactive or invalid")
		return
	}

	// whilst we are here - let's get the scopes and inject them into the request context
	// they may be useful later
	if irObj.Scope != nil {
		r.Header.Set("X-Tyk-Plugin-OAuth2Introspect-Scope", *irObj.Scope)
	}

	// strip the access token
	// r.Header.Del(authorizationHeaderName)
}

func accessTokenFromRequest(r *http.Request) string {
	auth := r.Header.Get(authorizationHeaderName)
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		// Nothing in Authorization header, try access_token
		// Empty string returned if there's no such parameter
		if err := r.ParseMultipartForm(1 << 20); err != nil && err != http.ErrNotMultipart {
			return ""
		}
		return r.Form.Get("access_token")
	}

	return split[1]
}

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

// writeUnauthorizedJSON writes a consistent JSON error body for 401 responses
func writeUnauthorizedJSON(w http.ResponseWriter, msg string) {
	type errBody struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(errBody{Error: "unauthorized", Message: msg})
}

// writeJSONError writes an error payload with custom status and code
func writeJSONError(w http.ResponseWriter, status int, code string, msg string) {
	type errBody struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errBody{Error: code, Message: msg})
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

func main() {}
