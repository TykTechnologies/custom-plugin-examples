import uuid

default_header_name = "Id"

def gen():
    id = uuid.uuid4()
    return str(id)