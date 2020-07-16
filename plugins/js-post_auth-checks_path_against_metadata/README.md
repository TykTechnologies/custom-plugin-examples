## What
A "post-key-auth" middleware that will check the key's metadata against the requested URI.
if the requested path contains a value in the metadata, the plugin will return a 401 error.

## How
Import the API definition included and place the JS code on the Gateway's file system to run it.