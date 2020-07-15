from tyk.decorators import *
from gateway import TykGateway as tyk

from rules import expressions
import re

@Hook
def BotDetectionMiddleware(request, session, spec):
    user_agent = request.get_header("User-Agent")
    if user_agent is None:
        return request, session

    match = None
    for expr in expressions:
        if re.match(expr, user_agent):
            match = True
            break

    if match:
        request.object.return_overrides.response_code = 401
        request.object.return_overrides.response_error = 'Not authorized'
    return request, session
