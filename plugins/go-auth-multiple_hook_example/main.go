package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/headers"
	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/user"
)

// called once plugin is loaded, this is where we put all initialization work for plugin
// i.e. setting exported functions, setting up connection pool to storage and etc.
func init() {
	var logger = log.Get()
	logger.Info("Processing Golang plugin init function version 1.2!!" )
	//Here you write the code for db connection
}


func ResponseSendCurrentTime(rw http.ResponseWriter, r *http.Request) {

	var logger = log.Get()
	apidef := ctx.GetDefinition(r)

	fmt.Println("Golang plugin - fmt example - API name is ", apidef.Name)

	logger.WithField("api-name", apidef.Name).Info("Processing HTTP request in Golang plugin!!" )

	//Demo injecting header to a request
	logger.WithField("api-name", apidef.Name).Info("Golang plugin - Adding header to a request before it goes upstream.")
	r.Header.Add("Foo", "Bar")

	logger.WithField("api-name", apidef.Name).Info("Golang plugin - ResponseSendCurrentTime")

	now := time.Now().String()


	getTime := r.URL.Query().Get("get_time")
	logger.WithField("api-name", apidef.Name).Info("Golang plugin - get_time is ", getTime)

	// check if we don't need to send reply
	if getTime != "yes" {
		// allow request to be processed and sent to upstream
		logger.WithField("api-name", apidef.Name).Info("Golang plugin - Adding current_time as a header in the request. Request to api will continue to the upstream")
		r.Header.Add("current_time", now)
		return
	}

	// send HTTP response from Golang plugin
	logger.WithField("api-name", apidef.Name).Info("Golang plugin - Setting the response header and body. Request will stop in this plugin and OK response will be returned")

	// prepare data to send
	replyData := map[string]string{
		"current_time": now,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	writeBody(rw, replyData)
}

func writeBody(rw http.ResponseWriter, replyJson map[string]string) error{

	var logger = log.Get()
	jsonData, err := json.Marshal(replyJson)
	if err != nil {
		logger.WithField("replyJson", replyJson).Error("Golang auth plugin: Failed to marshal")
		rw.WriteHeader(http.StatusInternalServerError)
		return err
	}

	rw.Write(jsonData)
	return nil
}


// ------------------------------------------------------------------
// Custom auth plugin code:

func getSessionByKey(key string) *user.SessionState {
	// here goes our logic to check if passed API key is valid and appropriate key session can be retrieved

	// perform auth (only one token "abc" is allowed)
	// Here you add code to query your database
	if key != "abc" {
		return nil
	}

	// return session
	return &user.SessionState{
		OrgID: "default",
		Alias: "abc-session",
	}
}

func MyPluginCustomAuthCheck(rw http.ResponseWriter, r *http.Request) {

	var logger = log.Get()
	apidef := ctx.GetDefinition(r)
	logger.WithField("api-name", apidef.Name).Info("Golang auth plugin - MyPluginCustomAuthCheck")

	// try to get session by API key
	key := r.Header.Get(headers.Authorization)
	session := getSessionByKey(key)
	if session == nil {
		// auth failed, reply with 403
		logger.WithField("api-name", apidef.Name).Info("Golang auth plugin - MyPluginCustomAuthCheck - failed")
		rw.WriteHeader(http.StatusForbidden)

		// prepare data to send
		replyData := map[string]string{
			"reason": "Access forbidden",
		}

		writeBody(rw, replyData)
		//jsonData, err := json.Marshal(replyData)
		//if err != nil {
		//	rw.WriteHeader(http.StatusInternalServerError)
		//
		//	rw.Write(jsonData)
		//	return
		//}

		return
	}

	logger.WithField("api-name", apidef.Name).Info("Golang auth plugin - MyPluginCustomAuthCheck - succeeded")


	// auth was successful, add session and key to request's context so other middle-wares can use it
	ctx.SetSession(r, session, key, true)
}

type OauthClientRequest struct {
	Method string `json:"method" msg:"method"`
	Body   string `json:"body" msg:"body"`
	Headers   string `json:"body" msg:"body"`
	Domain   string `json:"domain" msg:"domain"`
	Resource   string `json:"resource" msg:"resource"`

}
func main() {}


/*
func CustomAuthPluginOAuthOkta(rw http.ResponseWriter, r *http.Request)  () {

	var logger = log.Get()

	//Make api call to upstream target
	oauthClientRequest = OauthClientRequest {
		Method: "POST",
		Body: "{\"client_id\":\"{PASTE-YOUR-OWN-CLIEND-ID}\",\"client_secret\":\"{PASTE-YOUR-OWN-CLIEND-SECRET}\",\"audience\":\"auth0-id\",\"grant_type\":\"client_credentials\"}",
		Headers: "{\"content-type\": \"application/json\"}",
	  Domain: "https://{YOUR-ORG-NAME}.eu.auth0.com",
		Resource: "/oauth/token",
	};

	jsonData, err := json.Marshal(oauthClientRequest)
	if err != nil {
		logger.WithField("oauthClientRequest", oauthClientRequest).Error("Golang auth plugin: Failed to marshal")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(jsonData)


	var oauthClientRequestStr = JSON.stringify(oauthClientRequest)
	log("oauthClientRequest object: " + oauthClientRequestStr)

	rawlog("--- before get to upstream ---")
	oauthASResp = TykMakeHttpRequest(oauthClientRequestStr);
	rawlog("--- After get to upstream ---")

	log ('----')
	oauthASRespObj = JSON.parse(oauthASResp);
	var oauthASRespCode = JSON.parse(oauthASRespObj.Code);
	log('oauthASRespCode: ' + oauthASRespCode);

	var userRespCode = oauthASRespCode
	var userResponseBody = "empty body"
	if (oauthASRespCode != 200)
	{
		userResponseBody = "Error returned from AS (OAuth2.0 client credentials flow)."
		log("The request that was sent and failed to the AS: " + oauthClientRequestStr)
	}
	else
	{
		log('oauthASRespObj.Body: ' + oauthASRespObj.Body);
		oauthASRespBodyObj = JSON.parse(oauthASRespObj.Body)
		var backendReqAuthorization = oauthASRespBodyObj["access_token"]
		log ("backendReqAuthorization: " + backendReqAuthorization)

		backendRequest = {
		"Method": "GET",
		//"Body": "{\"empty\":\"body\"}",
		//"Headers": {"content-type":"application/json", "Authorization:": backendReqAuthorization},
			"Headers": {"Authorization": "Bearer " + backendReqAuthorization},
		"Domain": "http://0.0.0.0:80",
			"Resource": "/get"
	};
		var backendRequestStr = JSON.stringify(backendRequest)
		log('backendRequestStr: ' + backendRequestStr);
		var backendRequestObj = JSON.parse(backendRequestStr)

		rawlog("--- Before get to upstream ---")
		var backendResponse = TykMakeHttpRequest(backendRequestStr);
		rawlog("--- After get to upstream ---")

		backendRespObj = JSON.parse(backendResponse);
		userRespCode = JSON.parse(backendRespObj.Code);
		log('userRespCode: ' + userRespCode);

		if (userRespCode != 200)
		{
			userResponseBody = "Error returned from backend. request was:" + JSON.stringify(backendRequest)
		}
		else
		{
			backendRespBodyObj = JSON.parse(backendRespObj.Body)
			backendRespAuthorization = backendRespBodyObj.headers["Authorization"]
			userResponseBody = backendRespAuthorization
		}
	}

	var responseObject = {
	Body: "access_token from body resp of a backend: "+ userResponseBody,
		Headers: {
			"oauth-client": "client_credentials."
		},
	Code: userRespCode
	}
}*/

