package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/TykTechnologies/tyk/apidef/oas"
	tykctx "github.com/TykTechnologies/tyk/ctx"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// newTestOAS builds an OAS object with x-tyk-api-gateway extension for testing.
func newTestOAS(paths map[string]*openapi3.PathItem, listenPath string, pluginCfg map[string]interface{}) *oas.OAS {
	doc := &oas.OAS{}
	doc.OpenAPI = "3.0.3"
	doc.Info = &openapi3.Info{Title: "Test API", Version: "1.0"}

	doc.Paths = openapi3.NewPaths()
	for p, item := range paths {
		doc.Paths.Set(p, item)
	}

	xTyk := &oas.XTykAPIGateway{
		Info: oas.Info{
			ID:   "test-api",
			Name: "Test API",
		},
		Server: oas.Server{
			ListenPath: oas.ListenPath{
				Value: listenPath,
			},
		},
	}

	if pluginCfg != nil {
		xTyk.Middleware = &oas.Middleware{
			Global: &oas.Global{
				PluginConfig: &oas.PluginConfig{
					Data: &oas.PluginConfigData{
						Enabled: true,
						Value:   pluginCfg,
					},
				},
			},
		}
	}

	doc.Extensions = map[string]interface{}{
		"x-tyk-api-gateway": xTyk,
	}

	return doc
}

// mockRequest creates an httptest request with the OAS injected into context.
func mockRequest(method, path string, oasDef *oas.OAS, headers map[string]string) (*http.Request, *httptest.ResponseRecorder) {
	r := httptest.NewRequest(method, path, nil)
	if oasDef != nil {
		c := context.WithValue(r.Context(), tykctx.OASDefinition, oasDef)
		r = r.WithContext(c)
	}
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	rw := httptest.NewRecorder()
	return r, rw
}

// jsonResponse is a helper to build OAS response with inline JSON example.
func jsonResponse(code string, example interface{}) *openapi3.Response {
	return &openapi3.Response{
		Description: ptr("Test response"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: example,
			},
		},
	}
}

// xmlResponse builds an OAS response with inline XML string example.
func xmlResponse(code string, xmlBody string) *openapi3.Response {
	return &openapi3.Response{
		Description: ptr("XML response"),
		Content: openapi3.Content{
			"text/xml": &openapi3.MediaType{
				Example: xmlBody,
			},
		},
	}
}

// dualResponse builds a response with both JSON and XML content types.
func dualResponse(jsonExample interface{}, xmlExample string) *openapi3.Response {
	return &openapi3.Response{
		Description: ptr("Dual response"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: jsonExample,
			},
			"text/xml": &openapi3.MediaType{
				Example: xmlExample,
			},
		},
	}
}

// namedExamplesResponse builds a response with named examples.
func namedExamplesResponse(examples map[string]interface{}) *openapi3.Response {
	exRefs := make(openapi3.Examples)
	for name, value := range examples {
		exRefs[name] = &openapi3.ExampleRef{
			Value: &openapi3.Example{
				Value: value,
			},
		}
	}
	return &openapi3.Response{
		Description: ptr("Named examples"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Examples: exRefs,
			},
		},
	}
}

// schemaResponse builds a response with only a schema (no examples) for ExampleExtractor fallback.
func schemaResponse(schema *openapi3.Schema) *openapi3.Response {
	return &openapi3.Response{
		Description: ptr("Schema only"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{Value: schema},
			},
		},
	}
}

func ptr(s string) *string { return &s }

// ---------------------------------------------------------------------------
// Mock Detection Tests
// ---------------------------------------------------------------------------

func TestIsMockRequest_HeaderTrue(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set("X-Mock", "true")
	if !isMockRequest(r) {
		t.Error("expected mock request for X-Mock: true")
	}
}

func TestIsMockRequest_HeaderTRUE(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set("X-Mock", "TRUE")
	if !isMockRequest(r) {
		t.Error("expected mock request for X-Mock: TRUE (case insensitive)")
	}
}

func TestIsMockRequest_HeaderFalse(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set("X-Mock", "false")
	if isMockRequest(r) {
		t.Error("should not be mock request for X-Mock: false")
	}
}

func TestIsMockRequest_NoHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	if isMockRequest(r) {
		t.Error("should not be mock request without header")
	}
}

func TestIsMockRequest_SubdomainMockDot(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	r.Host = "mock.api.example.com"
	if !isMockRequest(r) {
		t.Error("expected mock request for Host: mock.api.example.com")
	}
}

func TestIsMockRequest_SubdomainMockDash(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	r.Host = "mock-api.example.com"
	if !isMockRequest(r) {
		t.Error("expected mock request for Host: mock-api.example.com")
	}
}

func TestIsMockRequest_SubdomainNormal(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	r.Host = "api.example.com"
	if isMockRequest(r) {
		t.Error("should not be mock request for Host: api.example.com")
	}
}

func TestIsMockRequest_SubdomainWithPort(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	r.Host = "mock.api.example.com:8080"
	if !isMockRequest(r) {
		t.Error("expected mock request for Host: mock.api.example.com:8080")
	}
}

// ---------------------------------------------------------------------------
// Config Parsing Tests
// ---------------------------------------------------------------------------

func TestParseMockConfig_Nil(t *testing.T) {
	cfg := parseMockConfig(nil)
	if !cfg.Enabled {
		t.Error("default should be enabled")
	}
	if cfg.DefaultContentType != "application/json" {
		t.Errorf("default content type = %q, want application/json", cfg.DefaultContentType)
	}
}

func TestParseMockConfig_EmptyMap(t *testing.T) {
	cfg := parseMockConfig(map[string]interface{}{})
	if !cfg.Enabled {
		t.Error("default should be enabled for empty map")
	}
}

func TestParseMockConfig_Disabled(t *testing.T) {
	cfg := parseMockConfig(map[string]interface{}{
		"mock_responder": map[string]interface{}{
			"enabled": false,
		},
	})
	if cfg.Enabled {
		t.Error("should be disabled when enabled: false")
	}
}

func TestParseMockConfig_CustomContentType(t *testing.T) {
	cfg := parseMockConfig(map[string]interface{}{
		"mock_responder": map[string]interface{}{
			"default_content_type": "text/xml",
		},
	})
	if cfg.DefaultContentType != "text/xml" {
		t.Errorf("content type = %q, want text/xml", cfg.DefaultContentType)
	}
}

func TestParseMockConfig_CORSOverride(t *testing.T) {
	cfg := parseMockConfig(map[string]interface{}{
		"mock_responder": map[string]interface{}{
			"cors": map[string]interface{}{
				"allow_origins": []interface{}{"https://portal.example.com"},
				"max_age":       float64(3600),
			},
		},
	})
	if len(cfg.CORS.AllowOrigins) != 1 || cfg.CORS.AllowOrigins[0] != "https://portal.example.com" {
		t.Errorf("CORS origins = %v, want [https://portal.example.com]", cfg.CORS.AllowOrigins)
	}
	if cfg.CORS.MaxAge != 3600 {
		t.Errorf("CORS max age = %d, want 3600", cfg.CORS.MaxAge)
	}
}

func TestParseMockConfig_Operations(t *testing.T) {
	cfg := parseMockConfig(map[string]interface{}{
		"mock_responder": map[string]interface{}{
			"operations": map[string]interface{}{
				"createLabel": map[string]interface{}{
					"preferred_content_type": "text/xml",
					"disabled":              true,
				},
			},
		},
	})
	opCfg, ok := cfg.Operations["createLabel"]
	if !ok {
		t.Fatal("expected createLabel in operations")
	}
	if opCfg.PreferredContentType != "text/xml" {
		t.Errorf("preferred CT = %q, want text/xml", opCfg.PreferredContentType)
	}
	if !opCfg.Disabled {
		t.Error("expected disabled = true")
	}
}

func TestParseMockConfig_MalformedSection(t *testing.T) {
	cfg := parseMockConfig(map[string]interface{}{
		"mock_responder": "not a map",
	})
	if !cfg.Enabled {
		t.Error("malformed section should fall back to defaults")
	}
}

// ---------------------------------------------------------------------------
// Map Helper Tests
// ---------------------------------------------------------------------------

func TestGetBoolFromMap_Float64(t *testing.T) {
	// JSON numbers are float64, but booleans should still work
	m := map[string]interface{}{"key": true}
	if !getBoolFromMap(m, "key", false) {
		t.Error("expected true")
	}
}

func TestGetIntFromMap_Float64(t *testing.T) {
	m := map[string]interface{}{"key": float64(42)}
	if got := getIntFromMap(m, "key"); got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestGetStringSliceFromMap_Mixed(t *testing.T) {
	m := map[string]interface{}{
		"key": []interface{}{"a", 42, "b"},
	}
	got := getStringSliceFromMap(m, "key")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("got %v, want [a b]", got)
	}
}

// ---------------------------------------------------------------------------
// Path Matching Tests
// ---------------------------------------------------------------------------

func TestMatchRequestToOperation_ExactPath(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/shipments": {
			Get: &openapi3.Operation{OperationID: "listShipments"},
		},
	}
	oasDef := newTestOAS(paths, "/api/v1/", nil)
	r, _ := mockRequest("GET", "/api/v1/shipments", oasDef, nil)

	matchedPath, _, op := matchRequestToOperation(r, oasDef)
	if op == nil {
		t.Fatal("expected match")
	}
	if matchedPath != "/shipments" {
		t.Errorf("matched = %q, want /shipments", matchedPath)
	}
	if op.OperationID != "listShipments" {
		t.Errorf("opID = %q, want listShipments", op.OperationID)
	}
}

func TestMatchRequestToOperation_PathParam(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/shipments/{shipmentId}": {
			Get: &openapi3.Operation{OperationID: "getShipment"},
		},
	}
	oasDef := newTestOAS(paths, "/api/v1/", nil)
	r, _ := mockRequest("GET", "/api/v1/shipments/SHP-001234", oasDef, nil)

	matchedPath, _, op := matchRequestToOperation(r, oasDef)
	if op == nil {
		t.Fatal("expected match for path parameter")
	}
	if matchedPath != "/shipments/{shipmentId}" {
		t.Errorf("matched = %q, want /shipments/{shipmentId}", matchedPath)
	}
}

func TestMatchRequestToOperation_NestedPathParam(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/shipments/{shipmentId}/tracking": {
			Get: &openapi3.Operation{OperationID: "getTracking"},
		},
	}
	oasDef := newTestOAS(paths, "/api/v1/", nil)
	r, _ := mockRequest("GET", "/api/v1/shipments/SHP-001234/tracking", oasDef, nil)

	_, _, op := matchRequestToOperation(r, oasDef)
	if op == nil {
		t.Fatal("expected match for nested path")
	}
	if op.OperationID != "getTracking" {
		t.Errorf("opID = %q, want getTracking", op.OperationID)
	}
}

func TestMatchRequestToOperation_NoMatch(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/shipments": {
			Get: &openapi3.Operation{OperationID: "listShipments"},
		},
	}
	oasDef := newTestOAS(paths, "/api/v1/", nil)
	r, _ := mockRequest("GET", "/api/v1/nonexistent", oasDef, nil)

	_, _, op := matchRequestToOperation(r, oasDef)
	if op != nil {
		t.Error("expected no match for undefined path")
	}
}

func TestMatchRequestToOperation_MethodMismatch(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/shipments": {
			Post: &openapi3.Operation{OperationID: "createShipment"},
		},
	}
	oasDef := newTestOAS(paths, "/api/v1/", nil)
	r, _ := mockRequest("GET", "/api/v1/shipments", oasDef, nil)

	_, _, op := matchRequestToOperation(r, oasDef)
	if op != nil {
		t.Error("expected no match for wrong method")
	}
}

func TestMatchRequestToOperation_RootPath(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/": {
			Get: &openapi3.Operation{OperationID: "root"},
		},
	}
	oasDef := newTestOAS(paths, "/api/", nil)
	r, _ := mockRequest("GET", "/api/", oasDef, nil)

	_, _, op := matchRequestToOperation(r, oasDef)
	// Root path "/" after stripping "/api" results in "/"
	if op == nil {
		t.Fatal("expected match for root path")
	}
}

func TestMatchRequestToOperation_VersionParam(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/service/rest/{version}/routing": {
			Post: &openapi3.Operation{OperationID: "determineRouting"},
		},
	}
	oasDef := newTestOAS(paths, "/route/", nil)
	r, _ := mockRequest("POST", "/route/service/rest/v3/routing", oasDef, nil)

	_, _, op := matchRequestToOperation(r, oasDef)
	if op == nil {
		t.Fatal("expected match for version path parameter")
	}
}

// ---------------------------------------------------------------------------
// Example Extraction Tests
// ---------------------------------------------------------------------------

func TestExtractExample_DirectInline(t *testing.T) {
	example := map[string]interface{}{"id": "SHP-001", "status": "in_transit"}

	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("OK"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: example,
			},
		},
	}})
	op := &openapi3.Operation{Responses: responses}

	code, ct, body, err := extractExample(op, 200, "application/json", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 200 {
		t.Errorf("code = %d, want 200", code)
	}
	if ct != "application/json" {
		t.Errorf("content type = %q, want application/json", ct)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if result["id"] != "SHP-001" {
		t.Errorf("id = %v, want SHP-001", result["id"])
	}
}

func TestExtractExample_NamedExampleSorted(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("401", &openapi3.ResponseRef{Value: namedExamplesResponse(map[string]interface{}{
		"zebra_format":  map[string]interface{}{"error": "zebra"},
		"apigee_format": map[string]interface{}{"error": "apigee"},
	})})

	op := &openapi3.Operation{Responses: responses}

	_, _, body, err := extractExample(op, 401, "application/json", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	// "apigee_format" comes first alphabetically
	if result["error"] != "apigee" {
		t.Errorf("expected first sorted example (apigee), got %v", result["error"])
	}
}

func TestExtractExample_NamedExampleByName(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("422", &openapi3.ResponseRef{Value: namedExamplesResponse(map[string]interface{}{
		"missing_field":  map[string]interface{}{"error": "missing"},
		"invalid_format": map[string]interface{}{"error": "invalid"},
	})})

	op := &openapi3.Operation{Responses: responses}

	_, _, body, err := extractExample(op, 422, "application/json", "missing_field")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["error"] != "missing" {
		t.Errorf("expected missing_field example, got %v", result["error"])
	}
}

func TestExtractExample_NamedExampleNotFound(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: namedExamplesResponse(map[string]interface{}{
		"default": map[string]interface{}{"ok": true},
	})})

	op := &openapi3.Operation{Responses: responses}

	_, _, _, err := extractExample(op, 200, "application/json", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent named example")
	}
}

func TestExtractExample_SchemaFallback(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: schemaResponse(&openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: openapi3.Schemas{
			"name":  &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
			"count": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"integer"}}},
		},
	})})

	op := &openapi3.Operation{Responses: responses}

	_, _, body, err := extractExample(op, 200, "application/json", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["name"] != "string" {
		t.Errorf("expected schema fallback name = 'string', got %v", result["name"])
	}
}

func TestExtractExample_StatusCodeNotFound(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("OK"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{Example: "ok"},
		},
	}})

	op := &openapi3.Operation{Responses: responses}

	_, _, _, err := extractExample(op, 404, "application/json", "")
	if err == nil {
		t.Error("expected error for missing status code 404")
	}
}

func TestExtractExample_204NoContent(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("204", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("No Content"),
	}})

	op := &openapi3.Operation{Responses: responses}

	code, _, body, err := extractExample(op, 204, "application/json", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 204 {
		t.Errorf("code = %d, want 204", code)
	}
	if len(body) != 0 {
		t.Errorf("expected empty body for 204, got %d bytes", len(body))
	}
}

func TestExtractExample_XMLStringPreserved(t *testing.T) {
	xmlBody := "<LabelResponse><id>LBL-001</id></LabelResponse>"
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("OK"),
		Content: openapi3.Content{
			"text/xml": &openapi3.MediaType{Example: xmlBody},
		},
	}})

	op := &openapi3.Operation{Responses: responses}

	_, ct, body, err := extractExample(op, 200, "text/xml", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != "text/xml" {
		t.Errorf("content type = %q, want text/xml", ct)
	}
	if string(body) != xmlBody {
		t.Errorf("XML body not preserved: got %q", string(body))
	}
}

// ---------------------------------------------------------------------------
// Content-Type Negotiation Tests
// ---------------------------------------------------------------------------

func TestResolveContentType_AcceptHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "text/xml")
	ct := resolveContentType(r, operationConfig{}, mockConfig{DefaultContentType: "application/json"})
	if ct != "text/xml" {
		t.Errorf("got %q, want text/xml", ct)
	}
}

func TestResolveContentType_AcceptWildcard(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "*/*")
	ct := resolveContentType(r, operationConfig{}, mockConfig{DefaultContentType: "application/json"})
	if ct != "application/json" {
		t.Errorf("got %q, want application/json (wildcard should use default)", ct)
	}
}

func TestResolveContentType_OperationOverride(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	ct := resolveContentType(r, operationConfig{PreferredContentType: "text/xml"}, mockConfig{DefaultContentType: "application/json"})
	if ct != "text/xml" {
		t.Errorf("got %q, want text/xml from operation config", ct)
	}
}

func TestResolveContentType_Default(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	ct := resolveContentType(r, operationConfig{}, mockConfig{DefaultContentType: "application/json"})
	if ct != "application/json" {
		t.Errorf("got %q, want application/json", ct)
	}
}

// ---------------------------------------------------------------------------
// Status Code Resolution Tests
// ---------------------------------------------------------------------------

func TestResolveStatusCode_Default200(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("OK")}})
	op := &openapi3.Operation{Responses: responses}
	r := httptest.NewRequest("GET", "/", nil)

	code, err := resolveStatusCode(op, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 200 {
		t.Errorf("code = %d, want 200", code)
	}
}

func TestResolveStatusCode_Prefers200Over201(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("201", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("Created")}})
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("OK")}})
	op := &openapi3.Operation{Responses: responses}
	r := httptest.NewRequest("GET", "/", nil)

	code, _ := resolveStatusCode(op, r)
	if code != 200 {
		t.Errorf("code = %d, want 200 (preferred over 201)", code)
	}
}

func TestResolveStatusCode_Falls201When200Missing(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("201", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("Created")}})
	op := &openapi3.Operation{Responses: responses}
	r := httptest.NewRequest("POST", "/", nil)

	code, _ := resolveStatusCode(op, r)
	if code != 201 {
		t.Errorf("code = %d, want 201", code)
	}
}

func TestResolveStatusCode_HeaderOverride(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("OK")}})
	op := &openapi3.Operation{Responses: responses}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Tyk-Accept-Example-Code", "404")

	code, err := resolveStatusCode(op, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 404 {
		t.Errorf("code = %d, want 404", code)
	}
}

func TestResolveStatusCode_InvalidHeader(t *testing.T) {
	op := &openapi3.Operation{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Tyk-Accept-Example-Code", "abc")

	_, err := resolveStatusCode(op, r)
	if err == nil {
		t.Error("expected error for non-numeric code")
	}
}

func TestValidateStatusCode_Injection(t *testing.T) {
	_, err := validateStatusCode("<script>alert(1)</script>")
	if err == nil {
		t.Error("expected error for XSS attempt in status code")
	}
}

func TestValidateStatusCode_OutOfRange(t *testing.T) {
	_, err := validateStatusCode("999")
	if err == nil {
		t.Error("expected error for code > 599")
	}
}

// ---------------------------------------------------------------------------
// CORS Tests
// ---------------------------------------------------------------------------

func TestSetCORSHeaders(t *testing.T) {
	rw := httptest.NewRecorder()
	cfg := corsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Authorization"},
		MaxAge:       3600,
	}
	setCORSHeaders(rw, cfg)

	if got := rw.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Allow-Origin = %q, want *", got)
	}
	if got := rw.Header().Get("Access-Control-Max-Age"); got != "3600" {
		t.Errorf("Max-Age = %q, want 3600", got)
	}
}

func TestHandleCORSPreflight_Options(t *testing.T) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("OPTIONS", "/test", nil)
	cfg := newDefaultConfig()

	handled := handleCORSPreflight(rw, r, cfg)
	if !handled {
		t.Error("expected OPTIONS to be handled")
	}
	if rw.Code != 204 {
		t.Errorf("status = %d, want 204", rw.Code)
	}
}

func TestHandleCORSPreflight_Get(t *testing.T) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	cfg := newDefaultConfig()

	handled := handleCORSPreflight(rw, r, cfg)
	if handled {
		t.Error("GET should not be handled as preflight")
	}
}

// ---------------------------------------------------------------------------
// Fail-Safe Behavior Tests
// ---------------------------------------------------------------------------

func TestMockResponder_NoMockHeader_PassThrough(t *testing.T) {
	r, rw := mockRequest("GET", "/test", nil, nil)
	MockResponder(rw, r)

	// Should not write anything — pass-through
	if rw.Code != 200 {
		// httptest.NewRecorder defaults to 200, verify no explicit write
		if rw.Body.Len() > 0 {
			t.Error("expected no body written for non-mock request")
		}
	}
}

func TestMockResponder_NilOAS_PassThrough(t *testing.T) {
	r, rw := mockRequest("GET", "/test", nil, map[string]string{"X-Mock": "true"})
	MockResponder(rw, r)

	if rw.Body.Len() > 0 {
		t.Error("expected no body written when OAS is nil")
	}
}

func TestMockResponder_DisabledConfig_PassThrough(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/test": {Get: &openapi3.Operation{OperationID: "test"}},
	}
	pluginCfg := map[string]interface{}{
		"mock_responder": map[string]interface{}{
			"enabled": false,
		},
	}
	oasDef := newTestOAS(paths, "/", pluginCfg)
	r, rw := mockRequest("GET", "/test", oasDef, map[string]string{"X-Mock": "true"})
	MockResponder(rw, r)

	if rw.Body.Len() > 0 {
		t.Error("expected no body written when config disabled")
	}
}

func TestMockResponder_UnmatchedPath_PassThrough(t *testing.T) {
	paths := map[string]*openapi3.PathItem{
		"/shipments": {Get: &openapi3.Operation{OperationID: "list"}},
	}
	oasDef := newTestOAS(paths, "/api/", nil)
	r, rw := mockRequest("GET", "/api/nonexistent", oasDef, map[string]string{"X-Mock": "true"})
	MockResponder(rw, r)

	if rw.Body.Len() > 0 {
		t.Error("expected no body written for unmatched path")
	}
}

// ---------------------------------------------------------------------------
// Analytics Suppression Test
// ---------------------------------------------------------------------------

func TestSuppressAnalytics(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	suppressAnalytics(r)

	val := r.Context().Value(tykctx.DoNotTrackThisEndpoint)
	if val != true {
		t.Error("expected DoNotTrackThisEndpoint = true after suppressAnalytics")
	}
}

// ---------------------------------------------------------------------------
// End-to-End Plugin Test
// ---------------------------------------------------------------------------

func TestMockResponder_EndToEnd_JSON200(t *testing.T) {
	example := map[string]interface{}{"id": "SHP-001", "status": "delivered"}

	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("OK"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{Example: example},
		},
	}})

	paths := map[string]*openapi3.PathItem{
		"/shipments": {
			Get: &openapi3.Operation{
				OperationID: "listShipments",
				Responses:   responses,
			},
		},
	}
	oasDef := newTestOAS(paths, "/api/v1/", nil)
	r, rw := mockRequest("GET", "/api/v1/shipments", oasDef, map[string]string{"X-Mock": "true"})

	MockResponder(rw, r)

	if rw.Code != 200 {
		t.Errorf("status = %d, want 200", rw.Code)
	}
	if rw.Header().Get("X-Mock-Response") != "true" {
		t.Error("missing X-Mock-Response header")
	}
	if rw.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rw.Header().Get("Content-Type"))
	}
	if rw.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("missing CORS headers")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rw.Body.Bytes(), &result); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if result["id"] != "SHP-001" {
		t.Errorf("id = %v, want SHP-001", result["id"])
	}

	// Verify analytics suppression
	val := r.Context().Value(tykctx.DoNotTrackThisEndpoint)
	if val != true {
		t.Error("analytics should be suppressed")
	}
}

func TestMockResponder_EndToEnd_XML(t *testing.T) {
	xmlBody := "<Label><id>LBL-001</id></Label>"

	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("OK"),
		Content: openapi3.Content{
			"text/xml":         &openapi3.MediaType{Example: xmlBody},
			"application/json": &openapi3.MediaType{Example: map[string]interface{}{"id": "LBL-001"}},
		},
	}})

	paths := map[string]*openapi3.PathItem{
		"/labels": {
			Post: &openapi3.Operation{
				OperationID: "createLabel",
				Responses:   responses,
			},
		},
	}
	oasDef := newTestOAS(paths, "/api/", nil)
	r, rw := mockRequest("POST", "/api/labels", oasDef, map[string]string{
		"X-Mock": "true",
		"Accept": "text/xml",
	})

	MockResponder(rw, r)

	if rw.Code != 200 {
		t.Errorf("status = %d, want 200", rw.Code)
	}
	if rw.Header().Get("Content-Type") != "text/xml" {
		t.Errorf("Content-Type = %q, want text/xml", rw.Header().Get("Content-Type"))
	}
	if string(rw.Body.Bytes()) != xmlBody {
		t.Errorf("body = %q, want %q", rw.Body.String(), xmlBody)
	}
}

func TestMockResponder_EndToEnd_Post201(t *testing.T) {
	example := map[string]interface{}{"id": "SHP-NEW", "status": "created"}

	responses := openapi3.NewResponses()
	responses.Set("201", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("Created"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{Example: example},
		},
	}})

	paths := map[string]*openapi3.PathItem{
		"/shipments": {
			Post: &openapi3.Operation{
				OperationID: "createShipment",
				Responses:   responses,
			},
		},
	}
	oasDef := newTestOAS(paths, "/api/", nil)
	r, rw := mockRequest("POST", "/api/shipments", oasDef, map[string]string{"X-Mock": "true"})

	MockResponder(rw, r)

	if rw.Code != 201 {
		t.Errorf("status = %d, want 201", rw.Code)
	}
}

func TestMockResponder_EndToEnd_DisabledOperation(t *testing.T) {
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr("OK"),
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{Example: "ok"},
		},
	}})

	paths := map[string]*openapi3.PathItem{
		"/disabled": {
			Get: &openapi3.Operation{
				OperationID: "disabledOp",
				Responses:   responses,
			},
		},
	}
	pluginCfg := map[string]interface{}{
		"mock_responder": map[string]interface{}{
			"operations": map[string]interface{}{
				"disabledOp": map[string]interface{}{
					"disabled": true,
				},
			},
		},
	}
	oasDef := newTestOAS(paths, "/", pluginCfg)
	r, rw := mockRequest("GET", "/disabled", oasDef, map[string]string{"X-Mock": "true"})

	MockResponder(rw, r)

	if rw.Body.Len() > 0 {
		t.Error("disabled operation should not produce a response")
	}
}
