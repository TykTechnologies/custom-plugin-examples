1. Create a file on the Tyk Gateway root directory under "middleware"

`middleware/injectHeader.js`

Here's our code for the JS file:
```
var testJSVMData = new TykJS.TykMiddleware.NewMiddleware({});

testJSVMData.NewProcessRequest(function(request, session, config) {

    console.log(JSON.stringify(request.Headers))
    request.SetHeaders['x-correlation-id'] = create_UUID();
	return testJSVMData.ReturnData(request, {});
});

function create_UUID(){
    var dt = new Date().getTime();
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (dt + Math.random()*16)%16 | 0;
        dt = Math.floor(dt/16);
        return (c=='x' ? r :(r&0x3|0x8)).toString(16);
    });
    return uuid;
}
```

2. In our API definition, we add:
```
"custom_middleware": {
      "pre": [
        {
          "name": "testJSVMData",
          "path": "./middleware/injectHeader.js",
          "require_session": false,
          "raw_body_only": false
        }
      ]
},
"driver": "otto"
```

3. Update the GW by pressing "Update" in the Dashboard or Restarting the Gateway if you are running in Headless

4. Test our API:
```
$ curl http://www.tyk-test.com:8080/api-test/get
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.64.1",
    "X-Amzn-Trace-Id": "Root=1-5ed7f235-a7362fdbc74ec4a6b2c6fcd6",
    "X-Correlation-Id": "b203bc34-5d67-444a-b359-d44c3db8029f"
  },
  "origin": "172.19.0.1, 206.174.179.211",
  "url": "http://httpbin.org/get"
}
```