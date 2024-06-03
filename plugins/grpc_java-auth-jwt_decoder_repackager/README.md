Java gRPC plugin
==

## Plugin overview

This repository provides a sample [gRPC](http://www.grpc.io/) plugin, written in Java, intended to work as part of [Tyk](https://tyk.io/). Gradle is used.

Simple auth logic is implemented, based on [Tyk custom middleware hooks](https://tyk.io/docs/tyk-api-gateway-v1-9/javascript-plugins/middleware-scripting/) logic.
A class implements the required hook methods.

## The hook

This plugin implements a single hook, it recieves a JWT, decodes it, and creates a new JWT using the Claims from the old JWT as well as adding some new ones.

## Running the gRPC server

	gradle run


## Enable GRPC in tyk.conf
```json
"coprocess_options": {
    "enable_coprocess": true,
    "coprocess_grpc_server": "tcp://localhost:5555"
  },
```


## Enable custom GRPC middleware in API definition
Add this to the Tyk API definition:

```json
"custom_middleware": {
      "pre": [],
      "post": [],
      "post_key_auth": [],
      "auth_check": {
        "name": "MyAuthHook",
        "path": "",
        "require_session": false,
        "raw_body_only": false
      },
      "response": [],
      "driver": "grpc",
      "id_extractor": {
        "extract_from": "",
        "extract_with": "",
        "extractor_config": {}
      }
    },
```
