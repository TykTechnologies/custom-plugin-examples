package main

import (
	"net/http"
	"strconv"

	"github.com/TykTechnologies/tyk/apidef/oas"
	tykctx "github.com/TykTechnologies/tyk/ctx"
)

// mockConfig holds the plugin configuration parsed from pluginConfig.data.value.mock_responder.
type mockConfig struct {
	Enabled            bool
	DefaultContentType string
	CORS               corsConfig
	Operations         map[string]operationConfig // keyed by operationId
}

// corsConfig controls CORS headers on mock responses.
type corsConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	MaxAge       int
}

// operationConfig allows per-operation overrides.
type operationConfig struct {
	PreferredContentType string
	PreferredExampleName string
	Disabled             bool
}

// newDefaultConfig returns sensible defaults when no plugin config is present.
func newDefaultConfig() mockConfig {
	return mockConfig{
		Enabled:            true,
		DefaultContentType: "application/json",
		CORS: corsConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders: []string{"*"},
			MaxAge:       86400,
		},
		Operations: make(map[string]operationConfig),
	}
}

// readPluginConfig extracts the pluginConfig.data.value map from the OAS API definition.
// Returns nil if the config is absent or disabled.
func readPluginConfig(r *http.Request) map[string]interface{} {
	oasDef := tykctx.GetOASDefinition(r)
	if oasDef == nil || oasDef.Extensions == nil {
		return nil
	}

	xTyk, ok := oasDef.Extensions["x-tyk-api-gateway"].(*oas.XTykAPIGateway)
	if !ok || xTyk == nil {
		return nil
	}
	if xTyk.Middleware == nil || xTyk.Middleware.Global == nil {
		return nil
	}

	pc := xTyk.Middleware.Global.PluginConfig
	if pc == nil || pc.Data == nil || !pc.Data.Enabled {
		return nil
	}

	return pc.Data.Value
}

// parseMockConfig parses the mock_responder section from the raw plugin config.
// Returns defaults if the section is missing or malformed.
func parseMockConfig(raw map[string]interface{}) mockConfig {
	cfg := newDefaultConfig()
	if raw == nil {
		return cfg
	}

	section, ok := raw["mock_responder"].(map[string]interface{})
	if !ok {
		return cfg
	}

	if enabled, exists := section["enabled"]; exists {
		cfg.Enabled = getBoolFromMap(section, "enabled", true)
		_ = enabled
	}

	if ct := getStringFromMap(section, "default_content_type"); ct != "" {
		cfg.DefaultContentType = ct
	}

	if corsSection, ok := section["cors"].(map[string]interface{}); ok {
		if origins := getStringSliceFromMap(corsSection, "allow_origins"); len(origins) > 0 {
			cfg.CORS.AllowOrigins = origins
		}
		if methods := getStringSliceFromMap(corsSection, "allow_methods"); len(methods) > 0 {
			cfg.CORS.AllowMethods = methods
		}
		if headers := getStringSliceFromMap(corsSection, "allow_headers"); len(headers) > 0 {
			cfg.CORS.AllowHeaders = headers
		}
		if maxAge := getIntFromMap(corsSection, "max_age"); maxAge > 0 {
			cfg.CORS.MaxAge = maxAge
		}
	}

	if opsSection, ok := section["operations"].(map[string]interface{}); ok {
		for opID, opRaw := range opsSection {
			opMap, ok := opRaw.(map[string]interface{})
			if !ok {
				continue
			}
			opCfg := operationConfig{
				PreferredContentType: getStringFromMap(opMap, "preferred_content_type"),
				PreferredExampleName: getStringFromMap(opMap, "preferred_example_name"),
				Disabled:             getBoolFromMap(opMap, "disabled", false),
			}
			cfg.Operations[opID] = opCfg
		}
	}

	return cfg
}

// getStringFromMap safely extracts a string from a map[string]interface{}.
func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getBoolFromMap safely extracts a bool from a map[string]interface{}.
// Returns defaultVal if the key is missing or not a bool.
func getBoolFromMap(m map[string]interface{}, key string, defaultVal bool) bool {
	v, ok := m[key]
	if !ok {
		return defaultVal
	}
	switch tv := v.(type) {
	case bool:
		return tv
	case string:
		b, err := strconv.ParseBool(tv)
		if err != nil {
			return defaultVal
		}
		return b
	default:
		return defaultVal
	}
}

// getIntFromMap safely extracts an int from a map[string]interface{}.
// Handles float64 (JSON default), int, and string representations.
func getIntFromMap(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch tv := v.(type) {
	case float64:
		return int(tv)
	case int:
		return tv
	case string:
		i, _ := strconv.Atoi(tv)
		return i
	default:
		return 0
	}
}

// getStringSliceFromMap safely extracts a []string from a map[string]interface{}.
func getStringSliceFromMap(m map[string]interface{}, key string) []string {
	v, ok := m[key]
	if !ok {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []string
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
