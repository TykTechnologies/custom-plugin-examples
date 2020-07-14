// ---- Sample middleware creation by end-user -----
var exampleJavaScriptMiddlewarePreHook = new TykJS.TykMiddleware.NewMiddleware({});

exampleJavaScriptMiddlewarePreHook.NewProcessRequest(function(request, session) {
    // You can log to Tyk console output by calling the built-in log() function:
    log("Hello from the Tyk JavaScript pre auth function")

    if (!request.Headers["Authorization"]) {
      request.ReturnOverrides.ResponseCode = 401
      request.ReturnOverrides.ResponseError = "Must Use Auth Token & Cert"
    }

    // You must return both the request and session metadata 
    return exampleJavaScriptMiddlewarePreHook.ReturnData(request, {});
});