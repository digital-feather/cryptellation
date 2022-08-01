# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: livetests.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0flivetests.proto\x12\tlivetests\"h\n\x07\x41\x63\x63ount\x12.\n\x06\x61ssets\x18\x01 \x03(\x0b\x32\x1e.livetests.Account.AssetsEntry\x1a-\n\x0b\x41ssetsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x02:\x02\x38\x01\"\x9e\x01\n\x15\x43reateLivetestRequest\x12@\n\x08\x61\x63\x63ounts\x18\x01 \x03(\x0b\x32..livetests.CreateLivetestRequest.AccountsEntry\x1a\x43\n\rAccountsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12!\n\x05value\x18\x02 \x01(\x0b\x32\x12.livetests.Account:\x02\x38\x01\"$\n\x16\x43reateLivetestResponse\x12\n\n\x02id\x18\x01 \x01(\x04\"Z\n SubscribeToLivetestEventsRequest\x12\n\n\x02id\x18\x01 \x01(\x04\x12\x15\n\rexchange_name\x18\x02 \x01(\t\x12\x13\n\x0bpair_symbol\x18\x03 \x01(\t\"#\n!SubscribeToLivetestEventsResponse\"\x16\n\x14LivetestEventRequest\"D\n\x15LivetestEventResponse\x12\x0c\n\x04type\x18\x01 \x01(\t\x12\x0c\n\x04time\x18\x02 \x01(\t\x12\x0f\n\x07\x63ontent\x18\x03 \x01(\t2\xc0\x02\n\x10LivetestsService\x12W\n\x0e\x43reateLivetest\x12 .livetests.CreateLivetestRequest\x1a!.livetests.CreateLivetestResponse\"\x00\x12x\n\x19SubscribeToLivetestEvents\x12+.livetests.SubscribeToLivetestEventsRequest\x1a,.livetests.SubscribeToLivetestEventsResponse\"\x00\x12Y\n\x0eListenLivetest\x12\x1f.livetests.LivetestEventRequest\x1a .livetests.LivetestEventResponse\"\x00(\x01\x30\x01\x42$Z\"/services/livetests/pkg/grpc/protob\x06proto3')



_ACCOUNT = DESCRIPTOR.message_types_by_name['Account']
_ACCOUNT_ASSETSENTRY = _ACCOUNT.nested_types_by_name['AssetsEntry']
_CREATELIVETESTREQUEST = DESCRIPTOR.message_types_by_name['CreateLivetestRequest']
_CREATELIVETESTREQUEST_ACCOUNTSENTRY = _CREATELIVETESTREQUEST.nested_types_by_name['AccountsEntry']
_CREATELIVETESTRESPONSE = DESCRIPTOR.message_types_by_name['CreateLivetestResponse']
_SUBSCRIBETOLIVETESTEVENTSREQUEST = DESCRIPTOR.message_types_by_name['SubscribeToLivetestEventsRequest']
_SUBSCRIBETOLIVETESTEVENTSRESPONSE = DESCRIPTOR.message_types_by_name['SubscribeToLivetestEventsResponse']
_LIVETESTEVENTREQUEST = DESCRIPTOR.message_types_by_name['LivetestEventRequest']
_LIVETESTEVENTRESPONSE = DESCRIPTOR.message_types_by_name['LivetestEventResponse']
Account = _reflection.GeneratedProtocolMessageType('Account', (_message.Message,), {

  'AssetsEntry' : _reflection.GeneratedProtocolMessageType('AssetsEntry', (_message.Message,), {
    'DESCRIPTOR' : _ACCOUNT_ASSETSENTRY,
    '__module__' : 'livetests_pb2'
    # @@protoc_insertion_point(class_scope:livetests.Account.AssetsEntry)
    })
  ,
  'DESCRIPTOR' : _ACCOUNT,
  '__module__' : 'livetests_pb2'
  # @@protoc_insertion_point(class_scope:livetests.Account)
  })
_sym_db.RegisterMessage(Account)
_sym_db.RegisterMessage(Account.AssetsEntry)

CreateLivetestRequest = _reflection.GeneratedProtocolMessageType('CreateLivetestRequest', (_message.Message,), {

  'AccountsEntry' : _reflection.GeneratedProtocolMessageType('AccountsEntry', (_message.Message,), {
    'DESCRIPTOR' : _CREATELIVETESTREQUEST_ACCOUNTSENTRY,
    '__module__' : 'livetests_pb2'
    # @@protoc_insertion_point(class_scope:livetests.CreateLivetestRequest.AccountsEntry)
    })
  ,
  'DESCRIPTOR' : _CREATELIVETESTREQUEST,
  '__module__' : 'livetests_pb2'
  # @@protoc_insertion_point(class_scope:livetests.CreateLivetestRequest)
  })
_sym_db.RegisterMessage(CreateLivetestRequest)
_sym_db.RegisterMessage(CreateLivetestRequest.AccountsEntry)

CreateLivetestResponse = _reflection.GeneratedProtocolMessageType('CreateLivetestResponse', (_message.Message,), {
  'DESCRIPTOR' : _CREATELIVETESTRESPONSE,
  '__module__' : 'livetests_pb2'
  # @@protoc_insertion_point(class_scope:livetests.CreateLivetestResponse)
  })
_sym_db.RegisterMessage(CreateLivetestResponse)

SubscribeToLivetestEventsRequest = _reflection.GeneratedProtocolMessageType('SubscribeToLivetestEventsRequest', (_message.Message,), {
  'DESCRIPTOR' : _SUBSCRIBETOLIVETESTEVENTSREQUEST,
  '__module__' : 'livetests_pb2'
  # @@protoc_insertion_point(class_scope:livetests.SubscribeToLivetestEventsRequest)
  })
_sym_db.RegisterMessage(SubscribeToLivetestEventsRequest)

SubscribeToLivetestEventsResponse = _reflection.GeneratedProtocolMessageType('SubscribeToLivetestEventsResponse', (_message.Message,), {
  'DESCRIPTOR' : _SUBSCRIBETOLIVETESTEVENTSRESPONSE,
  '__module__' : 'livetests_pb2'
  # @@protoc_insertion_point(class_scope:livetests.SubscribeToLivetestEventsResponse)
  })
_sym_db.RegisterMessage(SubscribeToLivetestEventsResponse)

LivetestEventRequest = _reflection.GeneratedProtocolMessageType('LivetestEventRequest', (_message.Message,), {
  'DESCRIPTOR' : _LIVETESTEVENTREQUEST,
  '__module__' : 'livetests_pb2'
  # @@protoc_insertion_point(class_scope:livetests.LivetestEventRequest)
  })
_sym_db.RegisterMessage(LivetestEventRequest)

LivetestEventResponse = _reflection.GeneratedProtocolMessageType('LivetestEventResponse', (_message.Message,), {
  'DESCRIPTOR' : _LIVETESTEVENTRESPONSE,
  '__module__' : 'livetests_pb2'
  # @@protoc_insertion_point(class_scope:livetests.LivetestEventResponse)
  })
_sym_db.RegisterMessage(LivetestEventResponse)

_LIVETESTSSERVICE = DESCRIPTOR.services_by_name['LivetestsService']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\"/services/livetests/pkg/grpc/proto'
  _ACCOUNT_ASSETSENTRY._options = None
  _ACCOUNT_ASSETSENTRY._serialized_options = b'8\001'
  _CREATELIVETESTREQUEST_ACCOUNTSENTRY._options = None
  _CREATELIVETESTREQUEST_ACCOUNTSENTRY._serialized_options = b'8\001'
  _ACCOUNT._serialized_start=30
  _ACCOUNT._serialized_end=134
  _ACCOUNT_ASSETSENTRY._serialized_start=89
  _ACCOUNT_ASSETSENTRY._serialized_end=134
  _CREATELIVETESTREQUEST._serialized_start=137
  _CREATELIVETESTREQUEST._serialized_end=295
  _CREATELIVETESTREQUEST_ACCOUNTSENTRY._serialized_start=228
  _CREATELIVETESTREQUEST_ACCOUNTSENTRY._serialized_end=295
  _CREATELIVETESTRESPONSE._serialized_start=297
  _CREATELIVETESTRESPONSE._serialized_end=333
  _SUBSCRIBETOLIVETESTEVENTSREQUEST._serialized_start=335
  _SUBSCRIBETOLIVETESTEVENTSREQUEST._serialized_end=425
  _SUBSCRIBETOLIVETESTEVENTSRESPONSE._serialized_start=427
  _SUBSCRIBETOLIVETESTEVENTSRESPONSE._serialized_end=462
  _LIVETESTEVENTREQUEST._serialized_start=464
  _LIVETESTEVENTREQUEST._serialized_end=486
  _LIVETESTEVENTRESPONSE._serialized_start=488
  _LIVETESTEVENTRESPONSE._serialized_end=556
  _LIVETESTSSERVICE._serialized_start=559
  _LIVETESTSSERVICE._serialized_end=879
# @@protoc_insertion_point(module_scope)
