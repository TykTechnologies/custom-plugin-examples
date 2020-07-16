// ---- Sample middleware creation by end-user -----
var exampleJavaScriptMiddlewarePostHook = new TykJS.TykMiddleware.NewMiddleware({});

exampleJavaScriptMiddlewarePostHook.NewProcessRequest(function(request, session, spec) {

  log('----------------------')
  log('session meta: ')
  log(JSON.stringify(session.meta_data)) // {"hello": "world"}

  var metadata_value = session.meta_data.hello // "world"
  
  // if the Request path contains "world", return an error
  if (request.RequestURI.indexOf(metadata_value) !== -1) {
    log('Found "world" - uh oh')
    request.ReturnOverrides.ResponseCode = 401
    request.ReturnOverrides.ResponseError = "Not allowed to access URI which intersects with metadata value \"world\""
  }

  // You must return both the request and session metadata    
  return exampleJavaScriptMiddlewarePostHook.ReturnData(request, session.meta_data);
});

// Log that middleware is initialised
log("JavaScript middleware is initialised");