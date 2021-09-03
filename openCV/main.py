import grpc
import logging
from grpc_basic import collect_pb2_grpc
from grpc_implementation.CollectorServicer import CollectorServicer
from concurrent import futures


def serve():
    """
    Start up a gRPC server.
    :return: None.
    """
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    collect_pb2_grpc.add_CollectorServicer_to_server(CollectorServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
