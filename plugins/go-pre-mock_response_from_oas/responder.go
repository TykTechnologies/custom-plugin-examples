package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/TykTechnologies/tyk/apidef/oas"
)

// findAndRespond matches the request to an OAS operation, extracts the example response,
// and writes it to the ResponseWriter. Returns true if a response was written.
func findAndRespond(rw http.ResponseWriter, r *http.Request, oasDef *oas.OAS, cfg mockConfig) bool {
	matchedPath, method, op := matchRequestToOperation(r, oasDef)
	if op == nil {
		logger.Debugf("[MockResponder] No OAS match for %s %s", r.Method, r.URL.Path)
		return false
	}

	// Check per-operation disable
	if op.OperationID != "" {
		if opCfg, ok := cfg.Operations[op.OperationID]; ok && opCfg.Disabled {
			logger.Debugf("[MockResponder] Operation %q disabled by config", op.OperationID)
			return false
		}
	}

	// Resolve content type
	opCfg := cfg.Operations[op.OperationID]
	contentType := resolveContentType(r, opCfg, cfg)

	// Resolve status code
	code, err := resolveStatusCode(op, r)
	if err != nil {
		logger.Debugf("[MockResponder] Status code error for %s %s: %v", method, matchedPath, err)
		return false
	}

	// Resolve example name
	exampleName := resolveExampleName(r, opCfg)

	// Extract example
	actualCode, actualCT, body, err := extractExample(op, code, contentType, exampleName)
	if err != nil {
		logger.Debugf("[MockResponder] No example for %s %s (code %d, ct %s): %v",
			method, matchedPath, code, contentType, err)
		return false
	}

	writeResponse(rw, actualCode, actualCT, body, matchedPath, cfg)
	logger.Debugf("[MockResponder] %s %s → %d %s (%d bytes)", method, matchedPath, actualCode, actualCT, len(body))
	return true
}

// matchRequestToOperation finds the OAS operation matching the request.
// Uses manual segment matching with path.Clean normalization (matching enforce-oas-auth pattern).
func matchRequestToOperation(r *http.Request, oasDef *oas.OAS) (string, string, *openapi3.Operation) {
	if oasDef.Paths == nil {
		return "", "", nil
	}

	reqPath := stripListenPath(r.URL.Path, oasDef)
	method := strings.ToUpper(r.Method)

	pathMap := oasDef.Paths.Map()

	// First: exact path match
	if pathItem, ok := pathMap[reqPath]; ok {
		if op := getOperationByMethod(pathItem, method); op != nil {
			return reqPath, method, op
		}
	}

	// Second: parameterized path match with scoring
	reqParts := splitPath(reqPath)
	var bestPath string
	var bestOp *openapi3.Operation
	bestScore := -1

	for oasPath, pathItem := range pathMap {
		oasParts := splitPath(oasPath)
		if len(oasParts) != len(reqParts) {
			continue
		}

		score := 0
		match := true
		for i, seg := range oasParts {
			if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
				score++
			} else if seg == reqParts[i] {
				score += 10
			} else {
				match = false
				break
			}
		}

		if match && score > bestScore {
			if op := getOperationByMethod(pathItem, method); op != nil {
				bestScore = score
				bestPath = oasPath
				bestOp = op
			}
		}
	}

	return bestPath, method, bestOp
}

// stripListenPath removes the API's listen path prefix and normalizes the path.
func stripListenPath(reqPath string, oasDef *oas.OAS) string {
	cleaned := path.Clean(reqPath)
	if oasDef.Extensions == nil {
		return cleaned
	}

	xTyk, ok := oasDef.Extensions["x-tyk-api-gateway"].(*oas.XTykAPIGateway)
	if !ok || xTyk == nil {
		return cleaned
	}

	lp := xTyk.Server.ListenPath.Value
	if lp == "" || lp == "/" {
		return cleaned
	}

	prefix := strings.TrimSuffix(lp, "/")
	if strings.HasPrefix(cleaned, prefix) {
		stripped := cleaned[len(prefix):]
		if stripped == "" {
			return "/"
		}
		return stripped
	}

	return cleaned
}

// splitPath splits a URL path into non-empty segments.
func splitPath(p string) []string {
	var parts []string
	for _, s := range strings.Split(p, "/") {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return parts
}

// getOperationByMethod returns the operation for a given HTTP method.
func getOperationByMethod(pathItem *openapi3.PathItem, method string) *openapi3.Operation {
	switch method {
	case "GET":
		return pathItem.Get
	case "POST":
		return pathItem.Post
	case "PUT":
		return pathItem.Put
	case "PATCH":
		return pathItem.Patch
	case "DELETE":
		return pathItem.Delete
	case "HEAD":
		return pathItem.Head
	case "OPTIONS":
		return pathItem.Options
	default:
		return nil
	}
}

// resolveContentType determines the response content type.
// Priority: Accept header > operation config > default config.
func resolveContentType(r *http.Request, opCfg operationConfig, cfg mockConfig) string {
	accept := r.Header.Get("Accept")
	if accept != "" && accept != "*/*" {
		return accept
	}
	if opCfg.PreferredContentType != "" {
		return opCfg.PreferredContentType
	}
	return cfg.DefaultContentType
}

// resolveStatusCode determines the desired HTTP status code for the mock response.
// Uses X-Tyk-Accept-Example-Code header if present, validates for safety.
func resolveStatusCode(op *openapi3.Operation, r *http.Request) (int, error) {
	headerCode := r.Header.Get("X-Tyk-Accept-Example-Code")
	if headerCode != "" {
		code, err := validateStatusCode(headerCode)
		if err != nil {
			return 0, err
		}
		return code, nil
	}

	// Auto-detect best status code from available responses
	if op.Responses == nil {
		return 200, nil
	}

	for _, code := range []int{200, 201, 202, 204} {
		if resp := op.Responses.Value(strconv.Itoa(code)); resp != nil {
			return code, nil
		}
	}

	return 200, nil
}

// validateStatusCode validates the X-Tyk-Accept-Example-Code header value.
// Prevents log injection by rejecting non-numeric input.
func validateStatusCode(s string) (int, error) {
	code, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid status code: not a number")
	}
	if code < 100 || code > 599 {
		return 0, fmt.Errorf("invalid status code: %d out of range 100-599", code)
	}
	return code, nil
}

// resolveExampleName determines which named example to use.
// Priority: X-Tyk-Accept-Example-Name header > operation config.
func resolveExampleName(r *http.Request, opCfg operationConfig) string {
	if name := r.Header.Get("X-Tyk-Accept-Example-Name"); name != "" {
		return name
	}
	return opCfg.PreferredExampleName
}

// extractExample finds and marshals the best example response from an OAS operation.
// Follows the exact priority from mw_mock_response.go:208-247:
//  1. Direct inline example (media.Example)
//  2. Named example by name (if specified)
//  3. First named example by sorted key (deterministic)
//  4. Schema-based fallback via oas.ExampleExtractor
func extractExample(op *openapi3.Operation, desiredCode int, contentType string, exampleName string) (int, string, []byte, error) {
	if op.Responses == nil {
		return 0, "", nil, fmt.Errorf("no responses defined")
	}

	// Find response for the desired code
	resp := op.Responses.Value(strconv.Itoa(desiredCode))
	if resp == nil || resp.Value == nil {
		return 0, "", nil, fmt.Errorf("no response definition for code %d", desiredCode)
	}

	// Find media type for content type
	media, actualCT := findMediaType(resp.Value.Content, contentType)
	if media == nil {
		// For 204 No Content, empty body is valid
		if desiredCode == 204 {
			return 204, contentType, nil, nil
		}
		return 0, "", nil, fmt.Errorf("no content type %q in response %d", contentType, desiredCode)
	}

	// Example extraction priority (matching gateway's mockFromOAS)
	var example interface{}

	// 1. Direct inline example
	if media.Example != nil {
		example = media.Example
	}

	// 2. Named examples
	if example == nil && len(media.Examples) > 0 {
		if exampleName != "" {
			exRef, ok := media.Examples[exampleName]
			if !ok || exRef == nil || exRef.Value == nil {
				return 0, "", nil, fmt.Errorf("named example %q not found", exampleName)
			}
			example = exRef.Value.Value
		} else {
			// Deterministic: first alphabetically sorted key
			keys := make([]string, 0, len(media.Examples))
			for k := range media.Examples {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				ex := media.Examples[k]
				if ex != nil && ex.Value != nil {
					example = ex.Value.Value
					break
				}
			}
		}
	}

	// 3. Schema-based fallback
	if example == nil && media.Schema != nil {
		example = oas.ExampleExtractor(media.Schema)
	}

	if example == nil {
		return 0, "", nil, fmt.Errorf("no example or schema available")
	}

	// Marshal based on content type
	body, err := marshalForContentType(example, actualCT)
	if err != nil {
		return 0, "", nil, fmt.Errorf("marshal error: %w", err)
	}

	return desiredCode, actualCT, body, nil
}

// findMediaType finds the best media type from the content map.
// Tries exact match first, then partial match (e.g., Accept: text/xml matches text/xml;charset=utf-8).
func findMediaType(content openapi3.Content, desired string) (*openapi3.MediaType, string) {
	if content == nil {
		return nil, ""
	}

	// Exact match
	if media := content.Get(desired); media != nil {
		return media, desired
	}

	// Partial match (Accept header might be "text/xml" while spec has "text/xml; charset=utf-8")
	desiredBase := strings.SplitN(desired, ";", 2)[0]
	desiredBase = strings.TrimSpace(desiredBase)
	for ct, media := range content {
		ctBase := strings.SplitN(ct, ";", 2)[0]
		ctBase = strings.TrimSpace(ctBase)
		if strings.EqualFold(ctBase, desiredBase) {
			return media, ct
		}
	}

	// If desired not found, try the first available
	for ct, media := range content {
		if media != nil {
			return media, ct
		}
	}

	return nil, ""
}

// marshalForContentType marshals an example value appropriate for the content type.
// XML string examples are returned as raw bytes. JSON uses json.Marshal.
func marshalForContentType(example interface{}, contentType string) ([]byte, error) {
	if example == nil {
		return nil, nil
	}

	isXML := strings.Contains(contentType, "xml")

	// XML string examples (like label endpoint) — return raw bytes
	if isXML {
		if s, ok := example.(string); ok {
			return []byte(s), nil
		}
	}

	// For JSON (or XML from schema fallback) — use json.Marshal
	// The gateway's mockFromOAS also uses json.Marshal for all content types.
	return json.Marshal(example)
}

// writeResponse writes the mock response with appropriate headers.
func writeResponse(rw http.ResponseWriter, statusCode int, contentType string, body []byte, matchedPath string, cfg mockConfig) {
	setCORSHeaders(rw, cfg.CORS)

	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("X-Mock-Response", "true")
	if matchedPath != "" {
		rw.Header().Set("X-Mock-Path", matchedPath)
	}

	rw.WriteHeader(statusCode)
	if len(body) > 0 {
		rw.Write(body)
	}
}
