import datetime
import base64
from tyk.decorators import *
from gateway import TykGateway as tyk
from ldap3 import Server, Connection, ALL


LDAP_SERVER='ipa.demo1.freeipa.org'


def prepareSessionState(session, username):
	# Set basic access rules for this API
	session.rate = 1000.0
	session.per = 1.0
	session.expires = -1

	# One hundred requests per 24 hours
	session.quota_max = 100
	session.quota_renews = int(datetime.datetime.now().timestamp()) + ((60 * 60) * 24)

	# Kill the cache entry in two minutes
	session.id_extractor_deadline = int(datetime.datetime.now().timestamp()) + 120

	return session 


def bindUser(username, password):
	cstr = 'uid={},cn=users,cn=accounts,dc=demo1,dc=freeipa,dc=org'.format(username)

	try:
		server = Server(LDAP_SERVER, get_info=ALL)
		conn = Connection(server, cstr, password, auto_bind=True)
	except:
		return False

	return True

def getUserAndPass(header):
	decoded = ""
	noLeader = header.replace("Basic ", "").replace("basic ", "")

	try:
		decoded = base64.b64decode(noLeader)
	except:
		return "", "", False 

	try:
		parts = str(decoded, "utf8").split(":")
	except:
		return ["", "", False] 

	if len(parts) != 2:
		return ["", "", False] 

	return [parts[0], parts[1], True]

@Hook
def MyLDAPMidleware(request, session, metadata, spec):
	auth_header = request.get_header('Authorization')
	upok = getUserAndPass(auth_header)

	username = upok[0]
	password = upok[1]
	ok = upok[2]
	
	if ok == False:
		return request, session, metadata

	if bindUser(username, password) == False:
		return request, session, metadata

	session = prepareSessionState(session, username)

	# Set meta data we can use later
	metadata["username"] = username 

	return request, session, metadata