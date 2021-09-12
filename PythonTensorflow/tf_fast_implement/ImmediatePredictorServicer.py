from tf_fast import tf_fast_pb2_grpc
from tf import tf_pb2
from model.predict import make_prediction
from model.model import create_model, compile_model
from static.const import WEIGHTS_PATH


class ImmediatePredictorServicer(tf_fast_pb2_grpc.ImmediatePredictorServicer):
    """
    The python gRPC Tensorflow server.
    """

    def __init__(self):
        """
        Constructor.
        """
        self.model = create_model()
        self.model.load_weights(WEIGHTS_PATH)
        compile_model(self.model)

    def ImmediatePred(self, request, context):
        """
        Receives the image info (name, path) from the server and returns a prediction immediately.
        :param request: tf_pb2.Image
        :param context: grpc._server._Context
        :return: tf_pb2.Prediction
        """
        name = request.name
        path = request.path
        b = make_prediction(self.model, path)
        return tf_pb2.Prediction(name=name, pred=b)
