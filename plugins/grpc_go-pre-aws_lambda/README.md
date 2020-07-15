# tyk-mw-grpcgo-lambda

## Start the gRPC plugin
```
go get -u github.com/TykTechnologies/tyk-mw-grpcgo-lambda
tyk-mw-grpcgo-lambda

---

2019/04/30 07:56:56 gRPC server listening on /tmp/foo.sock
2019/04/30 07:56:56 ### Available Functions ###
2019/04/30 07:56:56 --->  echo
2019/04/30 07:56:56 ###
2019/04/30 07:56:56 in your api definition, configure the custom middleware as follows, replacing FUNC_NAME with one of available functions printed above:
2019/04/30 07:56:56 "custom_middleware":{"pre":[{"name":"FUNC_NAME"}],"post":[],"post_key_auth":[],"auth_check":{"name":"","path":"","require_session":false},"response":[],"driver":"grpc","id_extractor":{"extract_from":"","extract_with":"","extractor_config":{}}}
```

Given that this is a demo, should you wish to customise listen address or region, please modify main.go and recompile.

https://github.com/TykTechnologies/tyk-mw-grpcgo-lambda/blob/master/main.go#L16-L20

## Configure Gateway to send messages to the gRPC plugin

Modify tyk.conf as follows:
```
"coprocess_options": {
  "enable_coprocess": true,
  "coprocess_grpc_server": "unix:///tmp/foo.sock",
},
```

## Configure your API definition to invoke the gRPC lambda plugin

```
"custom_middleware": {
  "pre": [
    {
      "name": "echo"
    }
  ],
  "driver": "grpc",
}
```

The above will attempt to invoke a lambda function called `echo` within `eu-west-2` (hardcoded in gRPC server main.go).
