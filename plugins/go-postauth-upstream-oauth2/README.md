# Upstream Oauth2.0 Plugin for Tyk Gateway

This plugin allows the Tyk Gateway to undergo a client_credentials flow for OAuth2.0 to obtain a JWT.
This subsequent JWT can be utilized for Authorizing access to APIs. 

## Configuration
Configuration is via environment variables.
```.env
TYK_CUSTOM_OAUTH_CLIENT_ID=YOUR_APPLICATION_CLIENT_ID
TYK_CUSTOM_OAUTH_CLIENT_SECRET=YOUR_APPLICATION_CLIENT_SECRET
TYK_CUSTOM_OAUTH_TOKEN_ENDPOINT=YOUR_APPLICATION_TOKEN_ENDPOINT
```

## Usage
Assuming you have an upstream that's protected by OAuth, this plugin will allow you to access that resource. 
Subsequently the returned Access Token will be cached in the Redis that is installed alongside Tyk-Gateway. 
You can layer on any Tyk native Authorization on-top of this plugin. 