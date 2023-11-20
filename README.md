# Custom Gateway Plugins

This is a repository that contains examples of [Tyk Plugins](https://tyk.io/docs/plugins/).  A Plugin is a custom middleware that is injected into the API request lifecycle, which further complements the built-in Tyk functionality such as authentication & rate limiting.

Here's the [different phases you can inject plugins](https://tyk.io/docs/concepts/middleware-execution-order/) in the request lifecycle.  A [response plugin](https://tyk.io/docs/plugins/response-plugins/) is also possible.

## Custom Go Plugin Examples

Language | Phase         | Description                                                 | Link 
-------- |---------------|-------------------------------------------------------------| --- 
Golang	| 	Pre	         | 	Injects client certificate attributes as a Header	         |	[Link](plugins/go-pre-cert_inject_dn)
GoLang	| 	Post-Auth	   | 	OAuth2 Introspection	                                      |	[Link](plugins/go-postauth-oauth2_introspection)
GoLang	| 	Post-Auth	   | 	Authorizes request against OPA	                            |	[Link](plugins/go-postauth-opa_integration)
Golang	| 	Post + Auth	 | 	Dummy one to test the 2 hooks in go	                       |	[Link](plugins/go-auth-multiple_hook_example)
Golang	| 	Pre	         | 	Checks Basic Auth creds  against an AWS DynamoDB instance	 | [Link](plugins/go-auth-basicauth_dynamodb)                  
Golang | Pre           | Custom Cache on upstream failure                            | [Link](https://gist.github.com/zalbiraw/d84ab1aef532ddc2b2ee3c6df81d836b)              
Golang | Pre           | Request funneling until cache is built                      | [Link](https://gist.github.com/zalbiraw/b1e25dfd2132cc55a05155f4ca291e19)
Golang | Pre           | Upstream URL rewrite based on header, query or body value   | [Link](plugins/go-header-url-rewrite-conditional)
Golang | Post          | Upstream OAuth2.0 (Client credentials flow)                 | [Link](plugins/go-postauth-upstream-oauth2)
Golang | Post          | Invoke AWS Lambda with IAM Credentials                      | [Link](plugins/go-postauth-invoke-aws-lambda)

### gRPC Plugin Languages

Language | Phase | Description                                                                                              | Link 
-------- | ----- |----------------------------------------------------------------------------------------------------------| --- 
gRPC (GoLang)	|	Pre	| 	Header Injection & Auth example	                                                                        |	[Link](plugins/grpc_go-auth-pre_headerinject_authhook) 
gRPC (GoLang)	|	Pre	| 	Invokes an AWS Lambda	                                                                                  |	[Link](plugins/grpc_go-pre-aws_lambda) 
gRPC (Node)	|	Pre / Autg	| 	Simple NODE example with access to config data Lambda	                                                                                  |	[Link](plugins/grpc_node-auth-simple) 
gRPC (Java)	|	Auth	| 	Decodes JWT, inserts a claim and resigns it	                                                            |	[Link](plugins/grpc_java-auth-jwt_decoder_repackager) 
gRPC (Java)	|	Post	| 	Inserts Metadata from the portal requested key as an HTTP header	                                       |	[Link](plugins/grpc_java-post-insert_metadata_as_header) 
gRPC (.NET)	|	Auth	| 	Performs auth check against a SQL server	                                                               |	[Link](plugins/grpc_dotnet-auth_sql_basicauth ) 
gRPC (Ruby)	|	Pre  | 	Modifies HTTP header	                                                                                   |	[Link](plugins/grpc_ruby-pre-header_modify) 
gRPC (Python)	|	Pre  | 	Inserts a HTTP header	                                                                                  |	[Link](plugins/grpc_python-pre-insert_header) 

### Javascript Plugin Languages

Language | Phase | Description                                                                                              | Link 
-------- | ----- |----------------------------------------------------------------------------------------------------------| --- 
Javascript	|	Pre	| 	Inserts tracing ID in header	                                                                           |	[Link](plugins/js-pre-insert_header)
Javascript	|	Pre	| 	Auth Token & mTLS protection	                                                                           |	[Link](plugins/js-pre-mtls_token_auth)
Javascript	|	Pre	| 	Evaluates the validity of a Tyk Token	                                                                  |	[Link](plugins/js-pre-token_inspection)
Javascript	|	Post	| 	Checks API requests against a WAF	                                                                      |	[Link](plugins/js-pre-post-waf)
Javascript	|	Post-Auth	| 	Checks the request path against the user's meta data.  If there is a cross-over, will deny the request	 |	[Link](plugins/js-post_auth-checks_path_against_metadata)

### Lua Plugin Languages

Language | Phase | Description                                                                                              | Link 
-------- | ----- |----------------------------------------------------------------------------------------------------------| --- 
Lua	|	Pre	| 	header injection	                                                                                       |	[Link](plugins/lua-pre-header_injection) 

### Python Plugin Languages

Language | Phase | Description                                                                                              | Link 
-------- | ----- |----------------------------------------------------------------------------------------------------------| --- 
Python	|	Auth	| 	Checks API requests against a hard-coded token	                                                         |	[Link](plugins/py-auth_example) 
Python	|	Auth	| 	Validates credentials against an LDAP server	                                                           |	[Link](plugins/py-auth-ldap_example) 
Python	|	Pre	| 	This plugin sends a message to a queue server, it uses kombu as the messaging library	                  |	[Link](plugins/py-pre-message_queue_kombo) 
Python	|	Pre	| 	This plugin sends log data to a Datadog agent.	                                                         |	[Link](plugins/py-pre-datadog_logger) 
Python	|	Pre	| 	This plugin sends log data to a Loggly HTTPS endpoint	                                                  |	[Link](plugins/py-pre-loggly_integration) 
Python	|	Pre	| 	This plugin will block requests from specific user agents, using regular expressions.	                  |	[Link](plugins/py-pre-bot_detection) 
Python	|	Pre  +  Post	| 	Inserts a correlation ID as a header	                                                                   |	[Link](plugins/py-pre_post-correlation_id_insert)
Python	|	Post	| 	Injects a signed JWT as Authorization Header	                                                           |	[Link](plugins/py-post-jwt-injection) 
Python  | Response | Modifies the header and body of a response                                                               | [Link](https://gist.github.com/oluwaseyeayinla/06605eff12e68c9920ccece1f545b4ac)


## Virtual Endpoints

[Virtual Endpoints](https://tyk.io/docs/advanced-configuration/compose-apis/virtual-endpoints/) are slightly different, more of a FaaS / Lambda as opposed to a plugin, and thus are treated differently

Language |  Description | Link 
-------- | ------------ | --- 
Javascript	|	Tyk as an OAuth2.0 Client in client_credentials flow in Auth0	|	https://gist.github.com/letzya/ba7c2cd833c11fac61ae4a1d1908f1dc
Javascript	|	Tyk as an OAuth2.0 Client in client_credentials flow in Azure	|	https://gist.github.com/letzya/7e852181643e871481a7997ae3d5b84a
Javascript	|	Demo body transform of response, XML to JSON using petstore's endpoint /pet/{id}	|	https://gist.github.com/letzya/7df4dbc37f2f075795995efb8e205d3e
Javascript	|	Make POST request with FormData to Upstream	|	[Link](plugins/ve_formdata-post)
Javascript	|	Create API Key via Dashboard API 	|	[Link](plugins/ve_createkey.md)

## Requests

Have a cool or useful idea to add to this list?  Feel free to open an issue.


## Developers

If adding an example, first off, thank you.

1.
Create a new directory the following name pattern:
```
<language>-<phase>-<description1>_<description2>_<description3>
```
For example:
```
js-pre-insert_header/
|- README.md
|- myplugin.js
|- apidef.js
```

2. Include a README with instructions, and the supporting files in the directory
