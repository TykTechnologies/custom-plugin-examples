Java gRPC plugin
==

## Plugin overview

This repository provides a sample [gRPC](http://www.grpc.io/) plugin, written in Java, intended to work as part of [Tyk](https://tyk.io/). Gradle is used.

## The hook

This plugin implements a single POST hook, it will inject metadata from a Portal Requested Key into the request as an HTTP Header

## Running the gRPC server

	gradle run


## Enable GRPC in tyk.conf
```json
"coprocess_options": {
    "enable_coprocess": true,
    "coprocess_grpc_server": "tcp://grpc_server_host:5555"
  },
```


## Enable custom GRPC middleware in API definition
Add this to the Tyk API definition:

```json
    "custom_middleware": {
      "pre": [],
      "post": [
        {
          "name": "HelloWorld",
          "path": "",
          "require_session": true,
          "raw_body_only": false
        }
      ],
      "post_key_auth": [],
      "auth_check": {
        "name": "",
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

## FAQ

I see this error in my gateway logs, why?

```
"time="Jul 23 18:27:32" level=error msg="Context variable type is not supported: map[string]interface {}
```

It's because we haven't set the context data.  That's ok - we're not using it for this use case and we can ignore this error.
