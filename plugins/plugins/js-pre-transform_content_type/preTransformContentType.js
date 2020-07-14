var preTransformContentType = new TykJS.TykMiddleware.NewMiddleware({});

preTransformContentType.NewProcessRequest(function(request, session, config) {
  log("Running sample PRE PROCESSOR preTransformContentType middleware");

  var thisSession = JSON.parse(TykGetKeyData(request.Headers["Authorization"], config.APIID))

  if (thisSession.status == "error" || thisSession.expires >= 1 && thisSession.expires < Math.round((new Date()).getTime() / 1000)) {
    log("Key expired, redirecting to login")
    request.ReturnOverrides.ResponseCode = 301
    request.ReturnOverrides.ResponseHeaders = {
       "Location": "http://anotherurl.com/"
    }
  }

  return preTransformContentType.ReturnData(request, session.meta_data);
});

log("Sample PRE middleware initialised");