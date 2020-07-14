# Custom Gateway Plugins

This is a repository that contains examples of [Tyk Plugins](https://tyk.io/docs/plugins/).  A Plugin is a custom middleware that is injected into the API request lifecycle, which further complements the built-in Tyk functionality such as authentication & rate limiting.

Here's the [different phases you can inject plugins](https://tyk.io/docs/concepts/middleware-execution-order/) in the request lifecycle.  A [response plugin](https://tyk.io/docs/plugins/response-plugins/) is also possible.

Whilst not technically a "plugin", [Virtual Endpoint](https://tyk.io/docs/advanced-configuration/compose-apis/virtual-endpoints/) examples are also included.

## Examples
Language | Phase | Description | Link 
-------- | ----- |------------ | --- 
Javascript	|	Pre	|	Inserts tracing ID in header	|	https://github.com/sedkis/tyk/blob/master/plugins/javascript/js-insert-header.md
Javascript	|	Pre	|	Auth Token & mTLS protection	|	https://github.com/sedkis/tyk/tree/master/plugins/javascript/auth-token-and-mtls-protection
Javascript	|	Pre	|	Evaluates the validity of a Tyk Token	|	https://github.com/sedkis/tyk/tree/master/plugins/javascript/pre-transform-content-type
Javascript	|	Post	|	Checks API requests against a WAF	|	https://github.com/sedkis/tyk/tree/master/plugins/javascript/waf
Golang	|	Post+Auth	|	Dummy one to test the 2 hooks in go	|	https://github.com/letzya/go-plugin-demo
gRPC (GoLang)	|	Pre	|	Header Injection	|	https://github.com/TykTechnologies/tyk-plugin-demo-golang
Python	|	Auth	|	Checks API requests against a hard-coded token	|	https://github.com/TykTechnologies/tyk-plugin-demo-python
Javascript	|	Virtual Endpoint	|	Tyk as an OAuth2.0 Client in client_credentials flow in Auth0	|	https://gist.github.com/letzya/ba7c2cd833c11fac61ae4a1d1908f1dc
Javascript	|	Virtual Endpoint	|	Tyk as an OAuth2.0 Client in client_credentials flow in Azure	|	https://gist.github.com/letzya/7e852181643e871481a7997ae3d5b84a
Javascript	|	Virtual Endpoint	|	Demo body transform of response, XML to JSON using petstore's endpoint /pet/{id}	|	https://gist.github.com/letzya/7df4dbc37f2f075795995efb8e205d3e
GoLang	|	Post-Auth	|	OAuth2 Introspection	|	https://github.com/asoorm/tyk-go-plugins/tree/master/oauth2_introspection
GoLang	|	Post-Auth	|	Authorizes request against OPA	|	https://github.com/asoorm/tyk-go-plugins/tree/master/authorize_opa

## Requests

Have a cool or useful idea to add to this list?  Feel free to open an issue.
