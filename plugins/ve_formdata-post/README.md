### Virtual Endpoint Example

This virtual endpoint will make an HTTP POST request to a mock upstream which echoes the request.
The request includes FormData.


## 1. Import the API into Tyk
Import the API definition located in this directory

## 2. Make API Call
Make an API call to invoke this serverless function:

```
$ curl localhost:8080/mock/virtual/

{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "hello": "world"
  },
  "headers": {
    "Accept-Encoding": "gzip",
    "Content-Length": "11",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "httpbin.org",
    "User-Agent": "Go-http-client/1.1",
    "X-Amzn-Trace-Id": "Root=1-5f2ae014-f2b1f58059f18f55b9982181"
  },
  "json": null,
  "origin": "99.242.138.81",
  "url": "http://httpbin.org/post"
}
```

Note the echo'd response contains the form parameters "hello=world"!

Here's the actual Javascript Virtual Endpoint code:
```javascript
function myUniqueFunctionName(request, session, config) {
  var responseObject = { 
    Body: "", 
    Code: 200 
  }
  
  response = logic(request)
  responseObject.Body = response.Body
  
  return TykJsResponse(responseObject, session.meta_data)
}

function logic(request) {
    //Make api call to upstream target
    postRequest = {
        "Method": "POST",
        "FormData": {"hello":"world"},
        "Headers": {"content-type":"application/x-www-form-urlencoded"},
        "Domain": "http://httpbin.org",
        "Resource": "/post"
    };
    var postRequestStr = JSON.stringify(postRequest)
    log("postRequestStr object: " + postRequestStr)
    
    rawlog("--- before upstream ---")
    responseObjectRaw = TykMakeHttpRequest(postRequestStr);
    rawlog("--- After upstream ---")
    
    responseObject = JSON.parse(responseObjectRaw);
    log('responseObject: ' + JSON.stringify(responseObject));
    return responseObject;
}
```
