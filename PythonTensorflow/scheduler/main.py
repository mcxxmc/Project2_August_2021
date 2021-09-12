from model.model import create_model, compile_model
# from test import test_model
from scheduler.Scheduler import Scheduler
from static.const import WEIGHTS_PATH

# CLIENT
model = create_model()
model.load_weights(WEIGHTS_PATH)
compile_model(model)

# launching the keras model by trying some tests
# test_model(model)

# create a scheduler thread
scheduler = Scheduler(model)
scheduler.run()
