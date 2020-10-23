### Virtual Endpoint Example

This virtual endpoint will create a Key through the Dashboard's API

```javascript
function myVirtualHandler(request, session, config) {
    //Make api call to upstream target
    newRequest = {
      "Headers": {"Authorization":"103cbe6f3c9443575cc52451cd051ecb"},
      "Method": "POST",
      "Body":keyrequest,
      "Domain": "http://www.tyk-test.com:3000",
      "resource": "/api/keys",
    };

    response = TykMakeHttpRequest(JSON.stringify(newRequest));
    log("response type: " + typeof response);
    log("response: " + response);
    
    var usableResponse = JSON.parse(response);
    var bodyObject = JSON.parse(usableResponse.Body);
  
    var responseObject = {
    Body: JSON.stringify(bodyObject.key_id),
    Headers: {
      "Set-Cookie": "session-id=" + JSON.stringify(bodyObject.key_id)
    },
    Code: usableResponse.Code
  }

  return TykJsResponse(responseObject, session.meta_data);
}

var keyrequest = '{"last_check":0,"certificate":null,"allowance":1000,"hmac_enabled":false,"hmac_string":"","basic_auth_data":{"password":""},"rate":1000,"per":60,"throttle_interval":-1,"throttle_retry_limit":-1,"expires":1563232857,"quota_max":-1,"quota_renews":1563146456,"quota_remaining":-1,"quota_renewal_rate":-1,"access_rights":{"06a1749bf11a4e2d7420224b5fb30158":{"api_id":"06a1749bf11a4e2d7420224b5fb30158","api_name":"service1","versions":["Default"],"allowed_urls":[],"limit":null}},"apply_policy_id":"","apply_policies":[],"tags":[],"jwt_data":{"secret":""},"meta_data":{},"alias":""}'
```
