# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import collect_pb2 as collect__pb2


class CollectorStub(object):
    """The greeting service definition.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.collectImage = channel.unary_unary(
                '/Collector/collectImage',
                request_serializer=collect__pb2.Empty.SerializeToString,
                response_deserializer=collect__pb2.ImageInfo.FromString,
                )


class CollectorServicer(object):
    """The greeting service definition.
    """

    def collectImage(self, request, context):
        """Receive the request to collect a new image using the camera;
        store the image in S3 and return the information of that image.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_CollectorServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'collectImage': grpc.unary_unary_rpc_method_handler(
                    servicer.collectImage,
                    request_deserializer=collect__pb2.Empty.FromString,
                    response_serializer=collect__pb2.ImageInfo.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'Collector', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Collector(object):
    """The greeting service definition.
    """

    @staticmethod
    def collectImage(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/Collector/collectImage',
            collect__pb2.Empty.SerializeToString,
            collect__pb2.ImageInfo.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
