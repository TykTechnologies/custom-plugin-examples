## What
A "post-key-auth" middleware that will check the key's metadata against the requested URI.
if the requested path contains a value in the metadata, the plugin will return a 401 error.

## How
Import the API definition included and place the JS code on the Gateway's file system to run it.

1. Create a key with metadata `"hello":"world"`

2. Curl it:

```
$ curl http://localhost:8080/javascript-middleware-api/get --header "Authorization: eyJv..."
{
  "hello": "from your API server!"
}


$ curl http://localhost:8080/javascript-middleware-api/world --header "Authorization: eyJv..."
{
    "error": "Not allowed to access URI which intersects with metadata value \"world\""
}
```
