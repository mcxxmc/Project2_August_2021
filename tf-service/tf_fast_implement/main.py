import grpc
from concurrent import futures
from static.const import GRPC_INSECURE_PORT
from tf_fast_implement.ImmediatePredictorServicer import ImmediatePredictorServicer
from tf_fast import tf_fast_pb2_grpc


def serve():
    """
    Starts a gPRC server.
    :return: None.
    """
    print("Server is started.")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    tf_fast_pb2_grpc.add_ImmediatePredictorServicer_to_server(ImmediatePredictorServicer(), server)
    server.add_insecure_port(GRPC_INSECURE_PORT)
    server.start()
    server.wait_for_termination()


# SERVER
serve()
