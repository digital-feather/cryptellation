# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from . import backtests_pb2 as backtests__pb2


class BacktestsServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CreateBacktest = channel.unary_unary(
                '/backtests.BacktestsService/CreateBacktest',
                request_serializer=backtests__pb2.CreateBacktestRequest.SerializeToString,
                response_deserializer=backtests__pb2.CreateBacktestResponse.FromString,
                )
        self.SubscribeToBacktestEvents = channel.unary_unary(
                '/backtests.BacktestsService/SubscribeToBacktestEvents',
                request_serializer=backtests__pb2.SubscribeToBacktestEventsRequest.SerializeToString,
                response_deserializer=backtests__pb2.SubscribeToBacktestEventsResponse.FromString,
                )
        self.ListenBacktest = channel.stream_stream(
                '/backtests.BacktestsService/ListenBacktest',
                request_serializer=backtests__pb2.BacktestEventRequest.SerializeToString,
                response_deserializer=backtests__pb2.BacktestEventResponse.FromString,
                )
        self.CreateBacktestOrder = channel.unary_unary(
                '/backtests.BacktestsService/CreateBacktestOrder',
                request_serializer=backtests__pb2.CreateBacktestOrderRequest.SerializeToString,
                response_deserializer=backtests__pb2.CreateBacktestOrderResponse.FromString,
                )
        self.Accounts = channel.unary_unary(
                '/backtests.BacktestsService/Accounts',
                request_serializer=backtests__pb2.AccountsRequest.SerializeToString,
                response_deserializer=backtests__pb2.AccountsResponse.FromString,
                )
        self.Orders = channel.unary_unary(
                '/backtests.BacktestsService/Orders',
                request_serializer=backtests__pb2.OrdersRequest.SerializeToString,
                response_deserializer=backtests__pb2.OrdersResponse.FromString,
                )


class BacktestsServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def CreateBacktest(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SubscribeToBacktestEvents(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListenBacktest(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def CreateBacktestOrder(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Accounts(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Orders(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_BacktestsServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CreateBacktest': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateBacktest,
                    request_deserializer=backtests__pb2.CreateBacktestRequest.FromString,
                    response_serializer=backtests__pb2.CreateBacktestResponse.SerializeToString,
            ),
            'SubscribeToBacktestEvents': grpc.unary_unary_rpc_method_handler(
                    servicer.SubscribeToBacktestEvents,
                    request_deserializer=backtests__pb2.SubscribeToBacktestEventsRequest.FromString,
                    response_serializer=backtests__pb2.SubscribeToBacktestEventsResponse.SerializeToString,
            ),
            'ListenBacktest': grpc.stream_stream_rpc_method_handler(
                    servicer.ListenBacktest,
                    request_deserializer=backtests__pb2.BacktestEventRequest.FromString,
                    response_serializer=backtests__pb2.BacktestEventResponse.SerializeToString,
            ),
            'CreateBacktestOrder': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateBacktestOrder,
                    request_deserializer=backtests__pb2.CreateBacktestOrderRequest.FromString,
                    response_serializer=backtests__pb2.CreateBacktestOrderResponse.SerializeToString,
            ),
            'Accounts': grpc.unary_unary_rpc_method_handler(
                    servicer.Accounts,
                    request_deserializer=backtests__pb2.AccountsRequest.FromString,
                    response_serializer=backtests__pb2.AccountsResponse.SerializeToString,
            ),
            'Orders': grpc.unary_unary_rpc_method_handler(
                    servicer.Orders,
                    request_deserializer=backtests__pb2.OrdersRequest.FromString,
                    response_serializer=backtests__pb2.OrdersResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'backtests.BacktestsService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class BacktestsService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def CreateBacktest(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/backtests.BacktestsService/CreateBacktest',
            backtests__pb2.CreateBacktestRequest.SerializeToString,
            backtests__pb2.CreateBacktestResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SubscribeToBacktestEvents(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/backtests.BacktestsService/SubscribeToBacktestEvents',
            backtests__pb2.SubscribeToBacktestEventsRequest.SerializeToString,
            backtests__pb2.SubscribeToBacktestEventsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListenBacktest(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(request_iterator, target, '/backtests.BacktestsService/ListenBacktest',
            backtests__pb2.BacktestEventRequest.SerializeToString,
            backtests__pb2.BacktestEventResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def CreateBacktestOrder(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/backtests.BacktestsService/CreateBacktestOrder',
            backtests__pb2.CreateBacktestOrderRequest.SerializeToString,
            backtests__pb2.CreateBacktestOrderResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Accounts(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/backtests.BacktestsService/Accounts',
            backtests__pb2.AccountsRequest.SerializeToString,
            backtests__pb2.AccountsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Orders(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/backtests.BacktestsService/Orders',
            backtests__pb2.OrdersRequest.SerializeToString,
            backtests__pb2.OrdersResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
