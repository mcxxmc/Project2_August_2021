import time

import grpc
from grpc_basic import tf_pb2_grpc, tf_pb2

from model.predict import make_prediction
import threading

from static.const import GRPC_GOLANG_INSECURE_PORT


class Scheduler(threading.Thread):
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
        threading.Thread.__init__(self)  # equals to super().__init__()
        self.model = model
        self.timeInterval = timeInterval
        self.t0 = None

        self.channel = grpc.insecure_channel(GRPC_GOLANG_INSECURE_PORT)
        self.stub = tf_pb2_grpc.CommunicatorStub(self.channel)

    def run(self) -> None:
        """
        Launch the scheduler.
        :return: None.
        """
        print("Scheduler launches.")

        self.t0 = time.time()

        while True:
            if int(time.time() - self.t0) == self.timeInterval:
                print("Job starts.")

                # Get all the images to predict
                response: tf_pb2.ImageArray = self.stub.RequestImages()
                namesPaths = {image.name: image.path for image in response.Images}
                r = []
                for name, path in namesPaths:
                    b = make_prediction(self.model, path)
                    r.append(tf_pb2.Prediction(name=name, pred=b))

                # send the predictions
                response: tf_pb2.Empty = self.stub.PostPredictions(Predictions=r)

                print("Job finishes. Timer resets.")
                self.t0 = time.time()
