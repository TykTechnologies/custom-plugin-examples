import json, requests

## Based on: https://www.loggly.com/docs/http-endpoint/
## Set the right URL for your API Key here:
loggly_url = 'https://logs-01.loggly.com/inputs/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/tag/{0}/'

def log(tag, record):
    requests.post(loggly_url.format(tag), json=record)
