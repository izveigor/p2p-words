# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import lemmatizer_pb2 as lemmatizer__pb2


class LemmatizersStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Lemmatize = channel.unary_unary(
                '/p2p.Lemmatizers/Lemmatize',
                request_serializer=lemmatizer__pb2.LemmatizerRequest.SerializeToString,
                response_deserializer=lemmatizer__pb2.LemmatizerResponse.FromString,
                )


class LemmatizersServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Lemmatize(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_LemmatizersServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Lemmatize': grpc.unary_unary_rpc_method_handler(
                    servicer.Lemmatize,
                    request_deserializer=lemmatizer__pb2.LemmatizerRequest.FromString,
                    response_serializer=lemmatizer__pb2.LemmatizerResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'p2p.Lemmatizers', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Lemmatizers(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Lemmatize(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/p2p.Lemmatizers/Lemmatize',
            lemmatizer__pb2.LemmatizerRequest.SerializeToString,
            lemmatizer__pb2.LemmatizerResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)