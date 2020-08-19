# go-dynamodb-basicauth-plugin

Built to be run natively as a plugin by Tyk Gateway.  

This is a "pre" plugin that will execute before Authentication is performed by Tyk.

This will just take the Client certificate in the API request and inject the cert issuer and fingerprint as headers to the upstream server.

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
