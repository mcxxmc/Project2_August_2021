import time
from model.predict import make_prediction


class Scheduler:

    def __init__(self, model):
        """
        Constructor.
        :param model: the keras model.
        """
        self.model = model
        self.t0 = None

    def launch(self, seconds: int = 30) -> None:
        """
        Launch the scheduler.
        :param seconds: int.
            The time between two jobs.
        :return: None.
        """
        print("Scheduler launches.")

        self.t0 = time.time()

        while True:
            if int(time.time() - self.t0) == seconds:
                print("Job starts.")
                make_prediction(self.model)
                print("Job finishes. Timer resets.")
                self.t0 = time.time()
