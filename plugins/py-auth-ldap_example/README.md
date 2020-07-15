# Tyk Python Plugin LDAP Auth Demo

This repository contains a demo LDAP authentication middleware for the Tyk API Gateway that will allow a user to authenticate against an API using their LDAP username and password as a Basic Auth header. 

The middleware makes use of Tyk's ID Extractor feature to limit round-trips to the LDAP server (caching credentials temporarily).

### Usage:

#### Clone the repo:

	git clone https://github.com/TykTechnologies/tyk-python-ldap-demo.git

#### Edit the `middleware.py` file 

to set the LDAP server to the correct host:

	LDAP_SERVER='ipa.demo1.freeipa.org'

If your LDAP server uses a different binding command, you will also need to edit the line that is:

	cstr = 'uid={},cn=users,cn=accounts,dc=demo1,dc=freeipa,dc=org'.format(username)

To match what you need.

#### Bundle everything up:

Create a bundle for Tyk to process using `tyk-cli`:

	tyk-cli bundle build -y -o ldap.zip

You will now need to place this somewhere where Tyk can fetch it, if in doubt, you can use the built-in python server:

    python -m SimpleHttpServer 8000

#### Ensure your tyk.conf has the correct settings:

	"enable_bundle_downloader": true,
	"bundle_base_url": "http://tyk-bundle-server:8000/bundles/", 
	"coprocess_options": {
		"enable_coprocess": true,
		"coprocess_grpc_server": ""
	},

You will need to be running the `tyk-python` version of the Tyk service.

#### Add the bundle to your API Definition:

To enable the bundle, add this to your API Definition, and make sure that you have it set up to use custom coprocess auth instead of another auth mechanism:

	"custom_middleware_bundle": "bundle.zip",
	"enable_coprocess_auth": true,

Once all that is ready, you can re-start Tyk.

You should now be able to make a request to the gateway like:

    http GET localhost:8181/testing/get Authorization:"Basic YWRtaW46U2VjcmV0MTIz"




