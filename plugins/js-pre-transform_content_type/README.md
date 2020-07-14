This is an example of a PRE middleware that will do the following:

    evaluate the validity of a token
    2a) If the token is valid, continue down the chain
    2b) if the token is expired or otherwise not found, the request will be terminated and a 301 redirect will be returned
