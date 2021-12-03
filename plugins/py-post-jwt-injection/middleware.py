from tyk.decorators import *
from gateway import TykGateway as tyk
from time import time
import os
import sys
import json

# Add vendor directory to module search path
parent_dir = os.path.abspath(os.path.dirname(__file__))
vendor_dir = os.path.join(parent_dir, 'vendor')
sys.path.append(vendor_dir)

import jwt
from cryptography.hazmat.primitives import serialization

@Hook
def PostMiddleware(request, session, spec):
    tyk.log("This is my post middleware", "info")

    customerID = request.get_header('CustomerID')
    #tyk.log("customerID: " + customerID, "info")
    
    payload = {
    "customerID": "1234",
    "userID": "andy@tyk.io",
    "country": "NL",
    "scope": "device:read"
    }

    my_secret = 'my_super_secret'
    private_key = open('/opt/tyk-gateway/id_rsa', 'r').read()
    key = serialization.load_ssh_private_key(private_key.encode(), password=b'')

    #token = jwt.encode(
    #    payload=payload,
    #    key=my_secret
    #)

    token = jwt.encode(
        payload=payload,
        key=key,
        algorithm='RS256'
    )

    tyk.log("jwt: " + token, "info")
    request.add_header("Authorization", token)

    return request, session
