# MockResponder — Tyk Go Pre-Plugin (OAS-driven mock responses)

A Go pre-plugin that returns OAS example responses for mock requests, so a developer portal's "Try It" flow (e.g. Stoplight Elements) can call your API without credentials or a live backend.

**One plugin, every gateway, decides per-request.** The plugin lives in the pre-middleware chain on all gateways. On each request it evaluates a trigger (mock subdomain or `X-Mock: true` header). If triggered, it pulls the OAS definition from the request context, matches the request to an operation, and returns the response example — short-circuiting auth, proxy, and response plugins. If not triggered, it returns immediately and the normal middleware chain proceeds (one string-compare overhead).

Works on **any** Tyk OAS API definition — no per-API plumbing. Response source priority: inline `example` → named `examples` → schema-synthesised via `oas.ExampleExtractor`.

Alternative to maintaining duplicate "mock twin" API definitions, which doubles the definition count at scale and bakes the routing decision into infrastructure.

## How It Works

```
Browser/curl                          Tyk Gateway
     │                                     │
     │  GET /shipments/v1/shipments        │
     │  X-Mock: true                       │
     │  ──────────────────────────────────>│
     │                                     │
     │            ┌─────────────────────┐  │
     │            │ MockResponder       │  │
     │            │ (Pre-Plugin)        │  │
     │            │                     │  │
     │            │ 1. X-Mock: true? ───│──│─── Yes ──┐
     │            │    or Host: mock.*  │  │          │
     │            │ 2. Read OAS Def    │  │          │
     │            │ 3. Match path      │  │          │
     │            │ 4. Extract example │  │          │
     │            │ 5. Write response  │  │          │
     │            └─────────────────────┘  │          │
     │                                     │          │
     │  200 OK + OAS Example               │  <───────┘
     │  X-Mock-Response: true              │
     │  <──────────────────────────────────│
     │                                     │
     │  (Auth, Proxy, Plugins all skipped) │
```

**Trigger:** `X-Mock: true` header OR `Host: mock.*` / `mock-*` subdomain.

**Fail-safe:** Without trigger, returns immediately — normal auth chain proceeds. Production traffic is unaffected (one string comparison overhead).

## Quick Start

### 1. Build

```bash
make build          # AMD64
make build-arm64    # ARM64 (Graviton/Fargate)
```

Output: `mockResponder_v5.11.0_linux_amd64.so` (~40 MB)

### 2. Deploy to Local Gateway

```bash
# Copy to middleware directory (volume-mounted into gateway container)
cp mockResponder_v5.11.0_linux_amd64.so ../../middleware/mockResponder.so

# Restart gateway to load plugin
docker compose restart tyk-gateway
```

### 3. Add to API Definition

Add `MockResponder` as **first** entry in `prePlugins` (must run before auth):

```json
"x-tyk-api-gateway": {
  "middleware": {
    "global": {
      "pluginConfig": {
        "driver": "goplugin"
      },
      "prePlugins": [
        {
          "enabled": true,
          "functionName": "MockResponder",
          "path": "/opt/tyk-gateway/middleware/mockResponder.so"
        }
      ]
    }
  }
}
```

### 4. Test

```bash
# Mock request (bypasses auth, returns OAS example)
curl http://localhost:8080/your-api/v1/endpoint -H "X-Mock: true"

# Normal request (auth chain as usual)
curl http://localhost:8080/your-api/v1/endpoint -H "Authorization: Bearer ..."

# Subdomain trigger (for portal Try It)
curl http://localhost:8080/your-api/v1/endpoint -H "Host: mock.api.example.com"

# Specific error code
curl http://localhost:8080/your-api/v1/endpoint \
  -H "X-Mock: true" \
  -H "X-Tyk-Accept-Example-Code: 401"

# Specific named example
curl http://localhost:8080/your-api/v1/endpoint \
  -H "X-Mock: true" \
  -H "X-Tyk-Accept-Example-Code: 422" \
  -H "X-Tyk-Accept-Example-Name: missing_field"

# XML content negotiation
curl http://localhost:8080/your-api/v1/labels/create \
  -X POST -H "X-Mock: true" -H "Accept: text/xml"
```

## Features

| Feature | Description |
|---------|-------------|
| **OAS Example Priority** | Matches gateway's built-in `mockResponse`: inline example > named example (alphabetical) > `ExampleExtractor` schema fallback |
| **Path Parameters** | `/shipments/{shipmentId}` matches `/shipments/SHP-001234` |
| **Content Negotiation** | Respects `Accept` header. JSON default, XML via `Accept: text/xml` |
| **Error Mocking** | `X-Tyk-Accept-Example-Code: 401` returns the 401 response example |
| **Named Examples** | `X-Tyk-Accept-Example-Name: apigee_format` selects a specific named example |
| **Schema Fallback** | APIs without examples get auto-generated responses via `oas.ExampleExtractor` |
| **CORS** | Configurable CORS headers. OPTIONS preflight handled automatically (required for portal) |
| **Analytics Suppression** | Sets `DoNotTrackThisEndpoint` — mock traffic does not pollute analytics |
| **Per-Operation Config** | Disable mock for specific operations, set preferred content type |
| **Fail-Safe** | No OAS, no match, or disabled → returns without writing, auth chain proceeds |
| **Panic Recovery** | Tyk's `defer recover()` catches panics — gateway never crashes |

## Configuration (Optional)

Add to `pluginConfig.data.value` in the API definition:

```json
"pluginConfig": {
  "driver": "goplugin",
  "data": {
    "enabled": true,
    "value": {
      "mock_responder": {
        "enabled": true,
        "default_content_type": "application/json",
        "cors": {
          "allow_origins": ["*"],
          "allow_methods": ["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"],
          "allow_headers": ["*"],
          "max_age": 86400
        },
        "operations": {
          "createLabel": {
            "preferred_content_type": "text/xml"
          },
          "internalEndpoint": {
            "disabled": true
          }
        }
      }
    }
  }
}
```

Without config, defaults are used (enabled, JSON, permissive CORS).

## Portal Integration (Stoplight Elements "Try It")

The plugin supports the portal's "Try It" flow:

1. **OAS `servers[0].url`** determines where Stoplight sends requests
2. **Subdomain trigger** works without custom headers: `servers[0].url = "https://mock.api.example.com/your-api/v1"`
3. **CORS preflight** is handled automatically — browsers send OPTIONS without custom headers, and the plugin responds with CORS headers even without `X-Mock`
4. **`tryItCredentialsPolicy="omit"`** in Stoplight Elements — no auth sent

### Portal Setup

Set the first OAS server to the mock subdomain:
```json
"servers": [
  {
    "url": "https://mock.api.example.com/shipments/v1",
    "description": "Mock (Try It)"
  }
]
```

The `mock.*` Host header triggers the plugin automatically.

## Response Headers

| Header | Value | Description |
|--------|-------|-------------|
| `X-Mock-Response` | `true` | Indicates this is a mock response |
| `X-Mock-Path` | `/shipments/{shipmentId}` | The matched OAS path template |
| `Access-Control-Allow-Origin` | `*` (configurable) | CORS header |
| `Content-Type` | `application/json` / `text/xml` | Resolved from Accept header or config |

## File Structure

```
go-plugin-mock-responder/
├── main.go            Entry point, mock detection, CORS preflight, analytics
├── config.go          Typed config struct, parsing from pluginConfig
├── responder.go       Path matching, example extraction, content negotiation
├── cors.go            CORS header writing, OPTIONS handling
├── main_test.go       ~50 unit tests (table-driven)
├── go.mod / go.sum    Dependencies (resolved by tyk-plugin-compiler)
├── Makefile           build / build-arm64 / test / clean
├── Dockerfile.gateway Multi-stage: test → compile → embed in gateway image
└── README.md          This file
```

## OAS Requirements

For best results, your OAS spec should have:

| Field | Required | Used For |
|-------|----------|----------|
| `operationId` | Recommended | Per-operation config, logging |
| `servers[0].url` | Required | Portal Try It target URL |
| `responses.{code}.content.{type}.example` | Recommended | Inline response examples |
| `responses.{code}.content.{type}.examples` | Optional | Named examples (error variants) |
| `components.schemas` | Optional | Fallback when no examples exist |

### Example Formats Supported

**Inline JSON:**
```json
"example": { "id": "SHP-001", "status": "delivered" }
```

**Named examples (for errors):**
```json
"examples": {
  "apigee_format": {
    "summary": "APIGEE error",
    "value": [{ "errorId": "APIGEE-403", "errorType": "ApiKeyDoesNotExistException" }]
  }
}
```

**XML string:**
```json
"text/xml": {
  "example": "<LabelResponse><id>LBL-001</id></LabelResponse>"
}
```

## Comparison: Plugin vs Mock Twin

| Aspect | MockResponder Plugin | Mock Twin (Generator) |
|--------|---------------------|----------------------|
| API Objects | 1 (production + plugin) | 2 (production + mock) |
| Config Change | Add prePlugin entry | Generate separate JSON |
| Auth on Production | Unchanged | Unchanged |
| CORS | Plugin handles it | In mock API definition |
| Analytics | `DoNotTrackThisEndpoint` | `trafficLogs.enabled: false` |
| Tyk Version | Recompile per version | Version-independent |
| Runtime Overhead | OAS clone per mock request | Zero (static config) |
| Portal Integration | Subdomain trigger | Separate listen path |

## Validated Test Results (v5.11.0)

| # | Test | Result |
|---|------|--------|
| 1 | GET JSON 200 inline example | PASS |
| 2 | GET path parameter matching | PASS |
| 3 | GET nested sub-resource | PASS |
| 4 | POST 201 created | PASS |
| 5 | POST XML content negotiation | PASS |
| 6 | POST JSON fallback for dual endpoint | PASS |
| 7 | Error code 401 (APIGEE format) | PASS |
| 8 | Error code 404 | PASS |
| 9 | Named example selection | PASS |
| 10 | OPTIONS CORS preflight | PASS |
| 11 | CORS Allow-Origin header | PASS |
| 12 | No X-Mock → auth chain (fail-safe) | PASS |
| 13 | Subdomain trigger via Host header | PASS |
| 14 | X-Mock-Response diagnostic header | PASS |
| 15 | X-Mock-Path diagnostic header | PASS |
| 16 | XML body content preserved | PASS |
| 17 | Performance: plugin <= mock twin latency | PASS |

## Known Limitations

1. **CORS preflight without X-Mock:** The plugin handles OPTIONS with `Access-Control-Request-Method` header even without `X-Mock`. This is necessary for portal integration but means OPTIONS on the production API will get CORS headers when the plugin is enabled. If this is undesirable, set `"mock_responder": { "enabled": false }` in the API config.

2. **Named example selection is deterministic (alphabetical):** Without `X-Tyk-Accept-Example-Name`, the first alphabetically-sorted named example is returned. This matches the gateway's built-in `mockResponse` behavior.

3. **Schema fallback returns JSON:** When using `oas.ExampleExtractor` for APIs without examples (e.g., XML-only Routing API), the fallback response is JSON-encoded. For realistic XML fallback, add inline XML examples to the OAS spec.

4. **Plugin compiler required:** Must recompile when upgrading Tyk Gateway version. Use `tykio/tyk-plugin-compiler:vX.Y.Z` matching the gateway version.
