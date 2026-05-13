#!/usr/bin/env bash
# MockResponder Plugin — Quick Setup for Testing
# Deploys the compiled plugin to the local gateway and syncs the test API.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
cd "${PROJECT_DIR}"

echo "=== MockResponder Plugin — Setup ==="
echo ""

# Step 1: Check for compiled .so
SO_FILE="${SCRIPT_DIR}/mockResponder_v5.11.0_linux_amd64.so"
if [ ! -f "${SO_FILE}" ]; then
    echo "ERROR: ${SO_FILE} not found."
    echo "Run 'make build' in ${SCRIPT_DIR} first."
    exit 1
fi
echo "1. Plugin binary found: $(ls -lh "${SO_FILE}" | awk '{print $5}')"

# Step 2: Copy to middleware directory
cp "${SO_FILE}" middleware/mockResponder.so
echo "2. Copied to middleware/mockResponder.so"

# Step 3: Ensure gateway loads from Dashboard
if grep -q '"use_db_app_configs": false' gateway/tyk.conf 2>/dev/null; then
    echo "   NOTE: gateway/tyk.conf has use_db_app_configs: false"
    echo "   The test API must be in apps/ OR set use_db_app_configs: true"
fi

# Step 4: Import test API to Dashboard (if Dashboard mode)
if [ -f .bootstrap-credentials ]; then
    source .bootstrap-credentials
    DASHBOARD_URL="http://localhost:${TYK_DASHBOARD_PORT:-3000}"

    echo "3. Importing test API to Dashboard..."
    RESP=$(curl -sf -X POST "${DASHBOARD_URL}/api/apis/oas" \
        -H "Authorization: ${TYK_USER_API_KEY}" \
        -H "Content-Type: application/json" \
        -d @apis/shipment-api-mock-plugin-test.json 2>&1) || true
    echo "   Response: $(echo "${RESP}" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('Status','?'), d.get('Message',''))" 2>/dev/null || echo "${RESP}")"
fi

# Step 5: Also copy to apps/ for file-based mode
cp apis/shipment-api-mock-plugin-test.json apps/ 2>/dev/null || true
echo "4. Copied test API to apps/ (file-based fallback)"

# Step 6: Restart gateway
echo "5. Restarting gateway..."
docker compose restart tyk-gateway 2>&1 | tail -1
sleep 5

# Step 7: Verify
echo ""
echo "=== Verification ==="

# Check API is loaded
API_COUNT=$(curl -s "http://localhost:8080/tyk/apis/" \
    -H "x-tyk-authorization: ${TYK_GW_SECRET:-352d20ee67be67f6340b4c0605b044b7}" 2>/dev/null | \
    python3 -c "import sys,json; print(len(json.load(sys.stdin)))" 2>/dev/null || echo "0")
echo "Gateway APIs loaded: ${API_COUNT}"

# Quick mock test
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    http://localhost:8080/shipments-mock-test/v1/shipments \
    -H "X-Mock: true" 2>/dev/null || echo "000")

if [ "${HTTP_CODE}" = "200" ]; then
    echo "Mock request: HTTP ${HTTP_CODE} — WORKING"
else
    echo "Mock request: HTTP ${HTTP_CODE} — NOT WORKING (check gateway logs)"
    echo "  Try: docker compose logs tyk-gateway --tail=20"
fi

echo ""
echo "=== Test Commands ==="
echo ""
echo "# JSON mock response:"
echo "curl http://localhost:8080/shipments-mock-test/v1/shipments -H 'X-Mock: true'"
echo ""
echo "# Path parameter:"
echo "curl http://localhost:8080/shipments-mock-test/v1/shipments/SHP-001234 -H 'X-Mock: true'"
echo ""
echo "# POST 201:"
echo "curl -X POST http://localhost:8080/shipments-mock-test/v1/shipments -H 'X-Mock: true'"
echo ""
echo "# XML response:"
echo "curl -X POST http://localhost:8080/shipments-mock-test/v1/labels/create -H 'X-Mock: true' -H 'Accept: text/xml'"
echo ""
echo "# Error code:"
echo "curl http://localhost:8080/shipments-mock-test/v1/shipments -H 'X-Mock: true' -H 'X-Tyk-Accept-Example-Code: 401'"
echo ""
echo "# Subdomain trigger:"
echo "curl http://localhost:8080/shipments-mock-test/v1/shipments -H 'Host: mock.api.example.com'"
echo ""
echo "# Normal request (should fail auth — plugin does NOT interfere):"
echo "curl http://localhost:8080/shipments-mock-test/v1/shipments"
