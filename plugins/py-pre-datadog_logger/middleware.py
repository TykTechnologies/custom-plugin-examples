from tyk.decorators import *
from gateway import TykGateway as tyk

from datadog import statsd

@Hook
def DatadogMiddleware(request, session, spec):
    tag = 'tyk.requests.api-{0}'.format(spec['APIID'])
    statsd.increment(tag)
    return request, session
