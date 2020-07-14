from tyk.decorators import *
from gateway import TykGateway as tyk

@Hook
def MyAuthMiddleware(request, session, metadata, spec):
    auth_header = request.get_header('Authorization')
    if auth_header == '47a0c79c427728b3df4af62b9228c8ae':
        session.rate = 1000.0
        session.per = 1.0
        metadata["token"] = "47a0c79c427728b3df4af62b9228c8ae"
    return request, session, metadata
