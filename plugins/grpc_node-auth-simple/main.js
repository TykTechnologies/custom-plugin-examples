const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

const packageDefinition = protoLoader.loadSync(
  path.resolve(__dirname, 'proto', 'coprocess_object.proto'),
  {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  }
);
const tyk = grpc.loadPackageDefinition(packageDefinition).coprocess;

const listenAddr = '127.0.0.1:5555',
    authHeader = 'Authorization'
    validToken = '71f6ac3385ce284152a64208521c592b'

// The dispatch function is called for every hook:
const dispatch = (call, callback) => {
  var obj = call.request
  // We dispatch the request based on the hook name, we pass obj.request which is the coprocess.Object:
  switch (obj.hook_name) {
    case 'MyPreMiddleware':
      preMiddleware(obj, callback)
      break
    case 'MyAuthMiddleware':
      authMiddleware(obj, callback)
      break
    default:
      callback(null, obj)
      break
  }
}

const preMiddleware = (obj, callback) => {
  // Config data
  console.log(JSON.stringify(obj.spec.config_data))

  var req = obj.request

  // req is the coprocess.MiniRequestObject, we inject a header using the "set_headers" field:
  req.set_headers = {
    'mycustomheader': 'mycustomvalue'
  }

  // Use this callback to finish the operation, sending back the modified object:
  callback(null, obj)
}

const authMiddleware = (obj, callback) => {
  var req = obj.request

  // We take the value from the "Authorization" header:
  var token = req.headers[authHeader]
  req.metadata

  // The token should be attached to the object metadata, this is used internally for key management:
  obj.metadata = {
    token: token
  }

  console.log(JSON.stringify(obj))

  // If the request token doesn't match the  "validToken" constant we return the call:
  if (token != validToken) {
    callback(null, obj)
    return
  }

  // At this point the token is valid and a session state object is initialized and attached to the coprocess.Object:
  var session = new tyk.SessionState()
  session.id_extractor_deadline = Date.now() + 100000000000
  obj.session = session
  callback(null, obj)
}

main = function() {
  server = new grpc.Server()
  server.addService(tyk.Dispatcher.service, {
      dispatch: dispatch,
      dispatchEvent: dispatch
  })
  server.bind(listenAddr, grpc.ServerCredentials.createInsecure())
  server.start()
}

main()