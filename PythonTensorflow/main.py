from model.model import create_model, compile_model
# from test import test_model
from scheduler.Scheduler import Scheduler


model = create_model()
model.load_weights("./weights/weight_2.h5")
compile_model(model)

# launching the keras model by trying some tests
# test_model(model)

scheduler = Scheduler(model)
scheduler.launch()
