package main

import (
	"net/http"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
)

var logger = log.Get()

// Injects key meta data into 
func InjectMetadata(rw http.ResponseWriter, r *http.Request) {
	session := ctx.GetSession(r)
	if session != nil {
		// Access session fields such as MetaData
		metaData := session.MetaData
		foo, ok := metaData["foo"].(string) // Type assert foo to string
		if !ok {
			// Handle the case where foo is not a string or foo does not exist
			logger.Error("Error: 'foo' is not a string or not found in metaData")
			return // or continue, depending on your error handling strategy
		}
		// Process metaData as needed
		r.Header.Add("X-Plugin-Inject", foo)
	}
}

func main() {}

func init() {
	logger.Info("--- Go custom plugin v4 init success! ---- ")
}
