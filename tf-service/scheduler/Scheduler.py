import logging
import time
import grpc
from tf import tf_pb2_grpc, tf_pb2
from model.predict import make_prediction
from static.const import GRPC_WEBSERVER_INSECURE_PORT


class Scheduler:
    """
    This class checks if there are new pictures to predict every 30 seconds (by default).
    Essentially a concurrent thread.
    """

    def __init__(self, model, timeInterval: int = 30):
        """
        Constructor.
        :param model: the keras model.
        :param timeInterval: int.
            The time between two jobs in seconds.
        """
        self.model = model
        self.timeInterval = timeInterval
        self.t0 = None

        self.channel = grpc.insecure_channel(GRPC_WEBSERVER_INSECURE_PORT)
        self.stub = tf_pb2_grpc.CommunicatorStub(self.channel)

    def run(self) -> None:
        """
        Launch the scheduler. Make gRPC requests as a client every 30 seconds.
        :return: None.
        """
        print("Scheduler launches.")

        self.t0 = time.time()

        while True:
            if int(time.time() - self.t0) == self.timeInterval:
                print("Job starts.")

                # gRPC client
                # Get all the images to predict
                try:
                    response: tf_pb2.ImageArray = self.stub.RequestImages(tf_pb2.TFStandard())
                    namesPaths = {image.name: image.path for image in response.Images}
                    r = []
                    for name, path in namesPaths.items():
                        b = make_prediction(self.model, path)
                        r.append(tf_pb2.Prediction(name=name, pred=b))

                    # send the predictions
                    self.stub.PostPredictions(tf_pb2.PredictionArray(Predictions=r))
                except Exception as e:
                    logging.error(e)

                print("Job finishes. Timer resets.")
                self.t0 = time.time()
