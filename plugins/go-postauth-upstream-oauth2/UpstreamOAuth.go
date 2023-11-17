package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/TykTechnologies/tyk/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/storage"
)

var logger = log.Get()
var store storage.RedisCluster
var requestInFlightTime = time.Duration(500) * time.Millisecond // Seconds value

var clientId string
var clientSecret string
var tokenEndpoint string
var grantType string
var tykCustomOAuthClientID = "TYK_CUSTOM_OAUTH_CLIENT_ID"
var tykCustomOAuthClientSecret = "TYK_CUSTOM_OAUTH_CLIENT_SECRET"
var tykCustomOAuthTokenEndpoint = "TYK_CUSTOM_OAUTH_TOKEN_ENDPOINT"

type AccessTokenStruct struct {
	Access_token       string `json:"access_token"`
	Expires_in         int64  `json:"expires_in"`
	Refresh_expires_in int    `json:"refresh_expires_in"`
	Refresh_token      string `json:"refresh_token"`
	Token_type         string `json:"token_type"`
	Not_before_policy  int    `json:"not-before-policy"`
	Session_state      string `json:"session_state"`
	Scope              string `json:"scope"`
}

func (m AccessTokenStruct) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

type OAuthError struct {
	Message string
}

func (e OAuthError) Error() string {
	return e.Message
}

func getStore(ctx context.Context) storage.RedisCluster {
	// Get Tyk Config
	conf := config.Global()

	// Create Redis Controller
	rc := storage.NewRedisController(ctx)
	logger.Debug("Created Redis Controller. Connected?", rc.Connected())

	store := storage.RedisCluster{KeyPrefix: "", HashKeys: conf.HashKeys, RedisController: rc}
	go rc.ConnectToRedis(ctx, nil, &conf)
	for i := 0; i < 5; i++ { // max 5 attempts - should only take 2
		if rc.Connected() {
			logger.Debug("Redis Controller connected")
			break
		}
		logger.Debug("Redis Controller not connected, will retry")

		time.Sleep(10 * time.Millisecond)
	}

	if !rc.Connected() {
		logger.Error("Could not connect to storage")
	}

	return store
}

func getCacheKey(clientId string) string {
	// Create base64 from clientId
	base64ClientId := base64.URLEncoding.EncodeToString([]byte(clientId))
	cacheKey := fmt.Sprintf("cache-%s", base64ClientId)
	logger.Debugf("cacheKey: %s", cacheKey)
	return cacheKey
}

func getAccessToken() (string, int64, error) {
	postForm := url.Values{}
	postForm.Add("client_id", clientId)
	postForm.Add("client_secret", clientSecret)
	postForm.Add("grant_type", grantType)
	requestBody := strings.NewReader(postForm.Encode())

	r, err := http.NewRequest("POST", tokenEndpoint, requestBody)
	if err != nil {
		logger.Errorf("Request creation error: %s", err.Error())
		return "", 0, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Exchange client credentials for an Access Token
	c := &http.Client{}
	res, err := c.Do(r)
	if err != nil {
		logger.Errorf("Error posting: %s", err.Error())
		return "", 0, err
	}

	if http.StatusOK != res.StatusCode {
		err = OAuthError{Message: fmt.Sprintf("Error OAuth request failed with %d", res.StatusCode)}
		logger.Error("Error OAuth request failed.")
		return "", 0, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Errorf("Error reading response from server: %s", err.Error())
		return "", 0, err
	}

	// Unmarshall the response into a AccessToken struct
	ats := AccessTokenStruct{}
	err = json.Unmarshal(bodyBytes, &ats)
	if nil != err {
		logger.Errorf("Error unmarshalling response body: %s", err.Error())
		return "", 0, err
	}

	return ats.Access_token, ats.Expires_in, nil
}

/*
UpstreamOAuth
Execute at Post Auth Middleware stage (assume initial JWT validation at this point)
Leverage ClientID & Secret that's from environment variables as a cacheKey
IF cacheKey exists && not expired, leverage it for a subsequent call to upstream (overwrite header)
ELSE make a post request to obtain a new JWT, update the cache key in Redis
*/
func UpstreamOAuth(rw http.ResponseWriter, r *http.Request) {
	// Create cacheKey from clientId, clientSecret and the request
	cacheKey := getCacheKey(clientId)

	// Using the cacheKey, check Redis for corresponding Access Token entry
	stringTTL, err := store.GetKeyTTL(cacheKey)
	if nil != err {
		logger.Errorf("Unable to check TTL for key: %s.", cacheKey)
		logger.Error(err.Error())
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	ttl := time.Duration(stringTTL) * time.Second

	// If cache exists and has more than requestInFlightTime seconds
	if ttl > requestInFlightTime {
		logger.Debug("Cache exists trying to retrieve token.")

		accessToken, err := store.GetKey(cacheKey)
		if nil != err {
			logger.Errorf("Error retrieving AccessToken from cache: %s.", err.Error())
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Set accessToken on request header
		r.Header.Set("Authorization", "Bearer "+accessToken)
		return
	}

	accessToken, expiresIn, err := getAccessToken()
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Set accessToken on request header
	r.Header.Set("Authorization", "Bearer "+accessToken)

	err = store.SetKey(cacheKey, accessToken, expiresIn)
	if err != nil {
		logger.Error("Unable to set AccessToken cache.")
		logger.Error(err.Error())
	}
}

func main() {}

func init() {
	clientId = os.Getenv(tykCustomOAuthClientID)
	clientSecret = os.Getenv(tykCustomOAuthClientSecret)
	tokenEndpoint = os.Getenv(tykCustomOAuthTokenEndpoint)
	grantType = "client_credentials"

	if clientId == "" {
		logger.Error(tykCustomOAuthClientID + " environment variable is not set")
		return
	}

	if clientSecret == "" {
		logger.Error(tykCustomOAuthClientSecret + " environment variable is not set")
		return
	}

	if tokenEndpoint == "" {
		logger.Error(tykCustomOAuthTokenEndpoint + " environment variable is not set")
		return
	}

	store = getStore(context.Background())

	logger.Info("--- Upstream Custom OAuth plugin initialized successfully! ---- ")
}
