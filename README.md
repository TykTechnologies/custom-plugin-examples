# Custom Gateway Plugins

This is a repository that contains examples of [Tyk Plugins](https://tyk.io/docs/plugins/).  A Plugin is a custom middleware that is injected into the API request lifecycle, which further complements the built-in Tyk functionality such as authentication & rate limiting.

Here's the [different phases you can inject plugins](https://tyk.io/docs/concepts/middleware-execution-order/) in the request lifecycle.  A [response plugin](https://tyk.io/docs/plugins/response-plugins/) is also possible.

Whilst not technically a "plugin", [Virtual Endpoint](https://tyk.io/docs/advanced-configuration/compose-apis/virtual-endpoints/) examples are also included.

## Examples
Language | Phase | Description | Link 
-------- | ----- |------------ | --- 
GoLang	|	Post-Auth	|	OAuth2 Introspection	|	[Link](plugins/go-postauth-oauth2_introspection)
GoLang	|	Post-Auth	|	Authorizes request against OPA	|	[Link](plugins/go-postauth-opa_integration)
Golang	|	Post + Auth	|	Dummy one to test the 2 hooks in go	|	[Link](plugins/go-auth-multiple_hook_example)
Golang	|	Auth	|	Checks Basic Auth creds  against an AWS DynamoDB instance	|	[Link](plugins/go-auth-basicauth_dynamodb)
gRPC (GoLang)	|	Pre	|	Header Injection & Auth example	|	[Link](plugins/grpc_go-auth-pre_headerinject_authhook) 
gRPC (Java)	|	Auth	|	Decodes JWT, inserts a claim and resigns it	|	[Link](plugins/grpc_java-auth-jwt_decoder_repackager) 
gRPC (.NET)	|	Auth	|	Performs auth check against a SQL server	|	[Link](plugins/grpc_dotnet-auth_sql_basicauth ) 
Javascript	|	Pre	|	Inserts tracing ID in header	|	[Link](plugins/js-pre-insert_header)
Javascript	|	Pre	|	Auth Token & mTLS protection	|	[Link](plugins/js-pre-mtls_token_auth)
Javascript	|	Pre	|	Evaluates the validity of a Tyk Token	|	[Link](plugins/js-pre-token_inspection)
Javascript	|	Post	|	Checks API requests against a WAF	|	[Link](plugins/js-pre-post-waf)
Javascript	|	Virtual Endpoint	|	Tyk as an OAuth2.0 Client in client_credentials flow in Auth0	|	https://gist.github.com/letzya/ba7c2cd833c11fac61ae4a1d1908f1dc
Javascript	|	Virtual Endpoint	|	Tyk as an OAuth2.0 Client in client_credentials flow in Azure	|	https://gist.github.com/letzya/7e852181643e871481a7997ae3d5b84a
Javascript	|	Virtual Endpoint	|	Demo body transform of response, XML to JSON using petstore's endpoint /pet/{id}	|	https://gist.github.com/letzya/7df4dbc37f2f075795995efb8e205d3e
Python	|	Auth	|	Checks API requests against a hard-coded token	|	[Link](plugins/py-auth_example) 
Python	|	Auth	|	Validates credentials against an LDAP server	|	[Link](plugins/py-auth-ldap_example) 
Python	|	Pre	|	This plugin sends a message to a queue server, it uses kombu as the messaging library	|	[Link](plugins/py-pre-message_queue_kombo) 
Python	|	Pre  +  Post	|	Inserts a correlation ID as a header	|	[Link](plugins/py-pre_post-correlation_id_insert) 

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