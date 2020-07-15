from tyk.decorators import *
from gateway import TykGateway as tyk

from correlation import gen, default_header_name

@Hook
def PreMiddleware(request, session, spec):
    id = gen()
    print("PRE, ID will be", id)
    request.add_header(default_header_name, id)
    return request, session

@Hook
def PostMiddleware(request, session, spec):
    id = request.get_header(default_header_name)
    print("POST, ID is", id)
     
    request.delete_header("Id")
    return request, session