package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TykTechnologies/tyk/header"
	"github.com/TykTechnologies/tyk/log"
)

var logger = log.Get()

// Your custom function to return the custom error message
func ReturnError(rw http.ResponseWriter, r *http.Request) {

	//your-custom-code-here

	//if header value != 'foo' return function below
	//if query param != 'x' return function below
	//if request is going to a specific target, return the error
	//based on our redirection rules, if request is redirected to x/y/z upstream, return the custom msg below

	errorHeader := r.Header.Get("errorHeader")
	logger.Info(fmt.Sprintf("ErrorHeader value is: %s", errorHeader))

	sendErrorResponseFromMiddleware(rw, errorHeader)
}

// Custom function that contains a list of custom error messages. 
func sendErrorResponseFromMiddleware(rw http.ResponseWriter, key string) {
	content := map[string]interface{}{
		"400": "This is a 400. This is also a custom error message returned from the plugin.",
		"401": "This is a 401. This is also a custom error message returned from the plugin.",
		"402": "This is a 402. This is also a custom error message returned from the plugin.",
		"403": "This is a 403. This is also a custom error message returned from the plugin.",
	}
	if val, ok := content[key]; ok {

		response := map[string]interface{}{key: val}

		data, err := json.Marshal(response)
		if err != nil {
			return
		}

		rw.Header().Set(header.ContentType, header.ApplicationJSON)
		errorCode, _ := strconv.Atoi(key)
		rw.WriteHeader(errorCode)
		rw.Write(data)
	}
}

func main() {}
func init() {
	logger.Info("--- Go custom plugin v4 init success! ---- ")
}
