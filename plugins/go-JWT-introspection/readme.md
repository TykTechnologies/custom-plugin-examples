## Introspection Plugin

This plugin implements introspection capabilities for JWT tokens.

### Notes
The client_id and client_secret are for the introspection client, not the client_id and client_secret of the application that is making the request.

Introspection should be enabled by the Identity Provider.
The introspection client should have the necessary permissions to introspect tokens.

### Compiling

Set the build parameters:
```bash
export BUILD_ID=12345
export GOOS=linux
export GOARCH=amd64
```

Then compile the plugin:
```bash
docker run --rm --name tyk-plugin-compiler \
  --network tyk \
  -v "$(pwd):/plugin-source" \
  -w /plugin-source \
  tykio/tyk-plugin-compiler-ee:v5.10.0 \
  plugin_introspect.so $BUILD_ID $GOOS $GOARCH
```

### Environment Variable Configuration

The plugin requires the following environment variables to be set on Tyk Gateway:

```bash
export OAUTH2_INTROSPECT_ENDPOINT="https://auth-server.com/oauth2/introspect"
export OAUTH2_INTROSPECT_CLIENT_ID="client_id_of_introspection_client"
export OAUTH2_INTROSPECT_CLIENT_SECRET="client_secret_of_introspection_client"
```
If using Tyk Helm Charts, use ExtraEnvs to set these environment variables.

### Suggested Execution Order

It is suggested to use this plugin as pre-authorization plugin with JWT middleware for scope-claim mapping to Tyk policies.



