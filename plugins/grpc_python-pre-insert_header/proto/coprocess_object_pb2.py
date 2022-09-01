# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: coprocess_object.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import coprocess_mini_request_object_pb2 as coprocess__mini__request__object__pb2
import coprocess_response_object_pb2 as coprocess__response__object__pb2
import coprocess_session_state_pb2 as coprocess__session__state__pb2
import coprocess_common_pb2 as coprocess__common__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x16\x63oprocess_object.proto\x12\tcoprocess\x1a#coprocess_mini_request_object.proto\x1a\x1f\x63oprocess_response_object.proto\x1a\x1d\x63oprocess_session_state.proto\x1a\x16\x63oprocess_common.proto\"\x85\x03\n\x06Object\x12&\n\thook_type\x18\x01 \x01(\x0e\x32\x13.coprocess.HookType\x12\x11\n\thook_name\x18\x02 \x01(\t\x12-\n\x07request\x18\x03 \x01(\x0b\x32\x1c.coprocess.MiniRequestObject\x12(\n\x07session\x18\x04 \x01(\x0b\x32\x17.coprocess.SessionState\x12\x31\n\x08metadata\x18\x05 \x03(\x0b\x32\x1f.coprocess.Object.MetadataEntry\x12)\n\x04spec\x18\x06 \x03(\x0b\x32\x1b.coprocess.Object.SpecEntry\x12+\n\x08response\x18\x07 \x01(\x0b\x32\x19.coprocess.ResponseObject\x1a/\n\rMetadataEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x1a+\n\tSpecEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"\x18\n\x05\x45vent\x12\x0f\n\x07payload\x18\x01 \x01(\t\"\x0c\n\nEventReply2|\n\nDispatcher\x12\x32\n\x08\x44ispatch\x12\x11.coprocess.Object\x1a\x11.coprocess.Object\"\x00\x12:\n\rDispatchEvent\x12\x10.coprocess.Event\x1a\x15.coprocess.EventReply\"\x00\x62\x06proto3')



_OBJECT = DESCRIPTOR.message_types_by_name['Object']
_OBJECT_METADATAENTRY = _OBJECT.nested_types_by_name['MetadataEntry']
_OBJECT_SPECENTRY = _OBJECT.nested_types_by_name['SpecEntry']
_EVENT = DESCRIPTOR.message_types_by_name['Event']
_EVENTREPLY = DESCRIPTOR.message_types_by_name['EventReply']
Object = _reflection.GeneratedProtocolMessageType('Object', (_message.Message,), {

  'MetadataEntry' : _reflection.GeneratedProtocolMessageType('MetadataEntry', (_message.Message,), {
    'DESCRIPTOR' : _OBJECT_METADATAENTRY,
    '__module__' : 'coprocess_object_pb2'
    # @@protoc_insertion_point(class_scope:coprocess.Object.MetadataEntry)
    })
  ,

  'SpecEntry' : _reflection.GeneratedProtocolMessageType('SpecEntry', (_message.Message,), {
    'DESCRIPTOR' : _OBJECT_SPECENTRY,
    '__module__' : 'coprocess_object_pb2'
    # @@protoc_insertion_point(class_scope:coprocess.Object.SpecEntry)
    })
  ,
  'DESCRIPTOR' : _OBJECT,
  '__module__' : 'coprocess_object_pb2'
  # @@protoc_insertion_point(class_scope:coprocess.Object)
  })
_sym_db.RegisterMessage(Object)
_sym_db.RegisterMessage(Object.MetadataEntry)
_sym_db.RegisterMessage(Object.SpecEntry)

Event = _reflection.GeneratedProtocolMessageType('Event', (_message.Message,), {
  'DESCRIPTOR' : _EVENT,
  '__module__' : 'coprocess_object_pb2'
  # @@protoc_insertion_point(class_scope:coprocess.Event)
  })
_sym_db.RegisterMessage(Event)

EventReply = _reflection.GeneratedProtocolMessageType('EventReply', (_message.Message,), {
  'DESCRIPTOR' : _EVENTREPLY,
  '__module__' : 'coprocess_object_pb2'
  # @@protoc_insertion_point(class_scope:coprocess.EventReply)
  })
_sym_db.RegisterMessage(EventReply)

_DISPATCHER = DESCRIPTOR.services_by_name['Dispatcher']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  _OBJECT_METADATAENTRY._options = None
  _OBJECT_METADATAENTRY._serialized_options = b'8\001'
  _OBJECT_SPECENTRY._options = None
  _OBJECT_SPECENTRY._serialized_options = b'8\001'
  _OBJECT._serialized_start=163
  _OBJECT._serialized_end=552
  _OBJECT_METADATAENTRY._serialized_start=460
  _OBJECT_METADATAENTRY._serialized_end=507
  _OBJECT_SPECENTRY._serialized_start=509
  _OBJECT_SPECENTRY._serialized_end=552
  _EVENT._serialized_start=554
  _EVENT._serialized_end=578
  _EVENTREPLY._serialized_start=580
  _EVENTREPLY._serialized_end=592
  _DISPATCHER._serialized_start=594
  _DISPATCHER._serialized_end=718
# @@protoc_insertion_point(module_scope)
