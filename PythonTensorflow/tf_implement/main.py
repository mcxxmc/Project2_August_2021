import grpc
from concurrent import futures
from static.const import GRPC_INSECURE_PORT
from tf import tf_pb2_grpc
from tf_implement.CommunicatorServicer import CommunicatorServicer


def serve():
    """
    Starts a gPRC server.
    :return: None.
    """
    print("Server is started.")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    tf_pb2_grpc.add_CommunicatorServicer_to_server(CommunicatorServicer(), server)
    server.add_insecure_port(GRPC_INSECURE_PORT)
    server.start()
    server.wait_for_termination()


# SERVER
serve()
