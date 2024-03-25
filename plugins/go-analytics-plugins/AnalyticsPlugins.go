package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TykTechnologies/tyk-pump/analytics"
	"github.com/TykTechnologies/tyk/log"
)

var logger = log.Get()

// Update raw record body content
func AnalyticsUpdateBody(record *analytics.AnalyticsRecord) {
	// Decode raw response string
	str, err := base64.StdEncoding.DecodeString(record.RawResponse)
	if err != nil {
		return
	}

	// Transform raw response string into http object
	var b = &bytes.Buffer{}
	b.Write(str)
	r := bufio.NewReader(b)

	var resp *http.Response
	resp, err = http.ReadResponse(r, nil)
	if err != nil {
		return
	}

	// Modify body content and update raw response value
	newBody := "Hello World!"
	resp.Body = ioutil.NopCloser(strings.NewReader(newBody))

	// Update related body contents
	resp.ContentLength = int64(len(newBody))

	// Set raw response value
	var bNew bytes.Buffer
	_ = resp.Write(&bNew)
	record.RawResponse = base64.StdEncoding.EncodeToString(bNew.Bytes())
}

// Delete raw response specified header
func AnalyticsDeleteHeader(record *analytics.AnalyticsRecord) {
	// Decode raw response string
	str, err := base64.StdEncoding.DecodeString(record.RawResponse)
	if err != nil {
		return
	}

	// Transform raw response string into http object
	var b = &bytes.Buffer{}
	b.Write(str)
	r := bufio.NewReader(b)

	var resp *http.Response
	resp, err = http.ReadResponse(r, nil)
	if err != nil {
		return
	}

	// Delete header and update raw response value
	resp.Header.Del("Server")

	// Set raw response value
	var bNew bytes.Buffer
	_ = resp.Write(&bNew)
	record.RawResponse = base64.StdEncoding.EncodeToString(bNew.Bytes())
}

// Delete full raw response
func AnalyticsDeleteRawResponse(record *analytics.AnalyticsRecord) {
	logger.Info("Raw response: ", record.RawResponse)

	// Decode raw response string
	str, err := base64.StdEncoding.DecodeString(record.RawResponse)
	if err != nil {
		logger.Info("Decode raw response error: ", str)
	}

	// Delete entire raw response string
	record.RawResponse = ""
}

// Delete full raw request
func AnalyticsDeleteRawRequest(record *analytics.AnalyticsRecord) {
	logger.Info("Raw request: ", record.RawRequest)

	// Decode raw request string
	str, err := base64.StdEncoding.DecodeString(record.RawRequest)
	if err != nil {
		logger.Info("Decode raw request error: ", str)
	}

	// Delete entire raw request string
	record.RawRequest = ""
}

func main() {}

func init() {
	logger.Info("--- Go custom plugin v4 init success! ---- ")