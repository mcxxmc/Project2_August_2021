import time
from model.predict import make_prediction
import threading


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
                make_prediction(self.model)
                print("Job finishes. Timer resets.")
                self.t0 = time.time()
