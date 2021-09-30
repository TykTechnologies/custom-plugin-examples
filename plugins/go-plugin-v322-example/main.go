package main

import (
	"github.com/TykTechnologies/tyk/log"
	"github.com/google/uuid"
	"net/http"
)

var logger = log.Get()


//Query param validation
func QueryParamValidation(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if len(query.Get("hello")) == 0 {
		logger.Warn("Go Plugin auth failure!")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	logger.Info("Go plugin auth passed.")
	r.Header.Set("X-GoPlugin-Random-Id", uuid.New().String())
}

// called once plugin is loaded, this is where we put all initialization work for plugin
// i.e. setting exported functions, setting up connection pool to storage and etc.
func init() {
	logger.Info("Initialising Example Go Plugin")
}

func main() {}
