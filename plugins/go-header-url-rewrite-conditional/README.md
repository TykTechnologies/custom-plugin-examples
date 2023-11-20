# Rewrite an upstream URL to include the value from a header, query parameter or payload conditional

This plugin allows the Tyk Gateway to rewrite the upstream URL of an incoming request if a header value, query parameter, or other value is present in the payload body sent in the request. If the conditional passes (storeNumber=xxxx), then include the value of the four digits in the upstream URL rewrite. For example, the incoming request is to localhost:8080/your-tyk-api/get?storeNumber=1221, the written URL is http://www.1221store.com.

## Configuration
The plugin is ready to run as is and mimics the native URL rewrite function inherent in the Gateway with extended capabilities. If you are looking to modify the logic, there are 3 helper functions (isValidStoreNumber, getXMLStoreNumber, setUpstream) that are used in the main function that executes the rewrite. Modifying the rewritten upstream URL can be done by changing line 74 in the code.

## Usage
Supposing you have the Tyk-Gateway deployed locally on port 8080, you will simply need to execute this plugin in the "pre" custom middleware section in your API definition. You should reference "RewriteUpstreamURL" as the function that is being implemented from the compiled Go code. In testing, you should see an upstream host lookup error as a response to the call, and notice that in the Gateway logs (or log browser) the proxied upstream will be the rewritten URL that is modified by the function (www.xxxxstore.com).
