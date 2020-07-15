from tyk.decorators import *
from gateway import TykGateway as tyk

from loggly import log

@Hook
def LogglyMiddleware(request, session, spec):
    tag = 'api-{0}'.format(spec['APIID'])
    record = {
        'user_agent': request.get_header('User-Agent')
    }
    log(tag, record)
    return request, session
