
# Go plugin demo

Demo Tyk go plugin for custom auth and post auth hooks.
Note the usage of Tyk packages.

## For Tyk version v2.9.3 mount pwd to "plugin-build"
docker run --rm -v `pwd`:/go/src/plugin-build tykio/tyk-plugin-compiler:v2.9.3 my-plugin-293.so

## For Tyk version v2.9.4.1 mount to "plugin-source"
docker run --rm -v `pwd`:/plugin-source tykio/tyk-plugin-compiler:2.9.4.1 my-go-plugin-2941.so

## Tesing the plugin: 

### Test the post hook plugin:

- `get_time=yes` - `current_time` is added as a header to your request and Tyk continues to reverse proxy the request to the backend.<br />
  `curl -s http://www.tyk-test.com:8080/go-plugin-demo/get?get_time=yes`
- `get_time != yes` - you get the `current_time` returned to you as the payload of the response. It's called "response ovverride". Tyk will stop the middleware execution chane and return response to the caller.<br />
`curl -s http://www.tyk-test.com:8080/go-plugin-demo/get?get_time=no`

### Test the custom auth hook plugin:
- Pass auth - authorization bearer is `abc`.<br />
 `curl -s 'http://www.localhost:8080/go-plugin-demo/get?get_time=yes' --header 'Authorization: abc' |jq .`
- Fail auth - authorization bearer is **not** `abc`.<br />
 `curl -s 'http://www.localhost:8080/go-plugin-demo/get?get_time=yes' --header 'Authorization: def' |jq .`

