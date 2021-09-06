import grpc
from concurrent import futures
from model.model import create_model, compile_model
# from test import test_model
from scheduler.Scheduler import Scheduler
from static.const import WEIGHTS_PATH, GRPC_INSECURE_PORT
from grpc_basic import tf_pb2_grpc
from grpc_implementation.CommunicatorServicer import CommunicatorServicer


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


# CLIENT
model = create_model()
model.load_weights(WEIGHTS_PATH)
compile_model(model)

# launching the keras model by trying some tests
# test_model(model)

# create a scheduler thread
scheduler = Scheduler(model)
scheduler.start()

# SERVER
serve()
