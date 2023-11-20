package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/regexp"
)

var logger = log.Get()
var storeRegex = `^\d{4}$`
var storeXMLRegex, _ = regexp.Compile(`storeNumber>(\d{4})<\/`)

func isValidStoreNumber(store string) bool {
	match, _ := regexp.MatchString(storeRegex, store)
	return match
}

func getXMLStoreNumber(body string) string {
	match := storeXMLRegex.FindStringSubmatch(body)

	if len(match) != 0 {
		return match[1]
	}

	return ""
}

func setUpstream(r *http.Request, upstream string) {
	apiDef := ctx.GetDefinition(r)
	apiDef.Proxy.EnableLoadBalancing = true
	apiDef.Proxy.StructuredTargetList = apidef.NewHostListFromList([]string{upstream})
	r.URL.Path = "store-tax-service-v1/"

}

func RewriteUpstreamURL(rw http.ResponseWriter, r *http.Request) {
	store := r.Header.Get("storeNumber")
	logger.Info(fmt.Sprintf("storeNumber header value: %s", store))

	if !isValidStoreNumber(store) {
		store = r.URL.Query().Get("storeNumber")
		logger.Info(fmt.Sprintf("storeNumber query parameter value: %s", store))

		if !isValidStoreNumber(store) {
			if r.Header.Get("Content-Type") == "application/xml" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Error(fmt.Sprintf("URL Rewrite custom plugin error: %s", err.Error()))
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				store = getXMLStoreNumber(string(body))
				logger.Info(fmt.Sprintf("storeNumber XML payload value: %s", store))

				if !isValidStoreNumber(store) {
					logger.Warn("No match for storeNumber in header, query parameter or XML")
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
			} else {
				logger.Warn("No match for storeNumber in header or query parameter")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	setUpstream(r, fmt.Sprintf("https://%sstore.com/", store))
}

func main() {}

func init() {
	logger.Info("--- Go custom plugin v4 init success! ---- ")
}
