# Certificate DN Injector

Built to be run natively as a GO plugin by Tyk Gateway.  

This is a "pre" plugin that will execute before authentication is performed by Tyk Gateway.

This will take the client certificate in the API request and inject the cert issuer and fingerprint as headers to the upstream server.

## Building From Docker
You can use the [Docker instructions](https://tyk.io/docs/plugins/supported-languages/golang/#building-a-golang-plugin) to generate a plugin binary.

In the plugin root, simply run:
`docker run --rm -v $(pwd):/plugin-source tykio/tyk-plugin-compiler:v3.0.0 mygeneratedfile.so`

This CLI will generate a file called `mygeneratefile.so`.  This is the file we need to mount onto the Gateway file system.  Skip down to "Setup your API" sections.

## Building From Source
### Generate the Binary file
Now we have to build the Go binary so that Tyk Gateway can natively run it.
`go build -o ./middleware/go/mygeneratedfile.so -buildmode=plugin`

### Build Tyk to run GOPLUGINS
Build Tyk with `goplugin` enabled.

In the root of tyk gateway, run:
`go build -tags 'goplugin' -o tyk .`

Then run the compiled Tyk
`./tyk`

# Setup the API
in API Designer, click on "Raw API Definition"
1. Set `"driver": "goplugin"`
2. Choose somewhere for your middleware to run in the cycle. ie:
```
"custom_middleware": {
      "pre": [
        {
          "name": "CertHeaderInject",
          "path": "./middleware/go/mygeneratedfile.so"
        }
      ],
   ```   
Pre is the phase in the cycle where it runs.
"name" has to be the name of the GO function
"path" is wherever you put the binary generated in step 1

## Example
With a mock upstream that echos HTTP requests:

```
$ curl https://tyk.gw/mtlsapi/get --cert clientpubkey.pem --key clientprivkey.pem
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.64.1",
    "X-Amzn-Trace-Id": "Root=1-5f3d6903-4b0afec0a05b7fc068087c00",
    "X-Client-Fingerprint": "53:2D:C3:73:5F:80:0A:E7:5C:1A:DD:E7:00:C7:4D:07:2D:9A:70:AA",
    "X-Client-Issuer": "my-self-signed-cert-issuer"
  },
  "origin": "172.23.0.1, 147.253.129.30",
  "url": "http://httpbin.org/get"
}
```
