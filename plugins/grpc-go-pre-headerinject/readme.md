# grpc-go-pre-headerinject

Allows a user to spin up a gRPC server in Go to service requests from Tyk.
Leverages the updated Protobuf definitions from Tyk. The included file `tykapidef_classic_grpc_go_header_injection.json` can be imported into Tyk and leverages the plugin.

## Quickstart
```zsh
# Navigate into the repo
cd grpc-go-pre-headerinject

go mod tidy

go run main.go
```

The gRPC server will now be running on localhost:5555. Import the file `tykapidef_classic_grpc_go_header_injection.json` into your Tyk installation and then invoke the resultant API to see the header injected.

## Note
The go.mod is sourced from [here](https://github.com/TykTechnologies/tyk/blob/master/go.mod).

The flatmap error can be overcome with the following replace statement in the go.mod file:

`replace github.com/hashicorp/terraform v1.0.1 => github.com/hashicorp/terraform v0.14.11`