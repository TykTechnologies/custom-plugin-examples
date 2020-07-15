from tyk.decorators import *
from gateway import TykGateway as tyk

import plugin as queue

queue.setup()

@Hook
def QueueMiddleware(request, session, spec):
    user_agent = request.get_header('User-Agent')
    queue.put(user_agent)
    return request, session
