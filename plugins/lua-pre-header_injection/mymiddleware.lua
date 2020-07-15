function MyPreMiddleware(request, session, spec)
  tyk.req.set_header("customheader", "customvalue")
  return request, session
end
