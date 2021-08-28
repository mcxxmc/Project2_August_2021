from common import create_model, load_image, compile_model
from common import S3_TO_PREDICT
from test import test_model
import os


model = create_model()
model.load_weights("./weights/weight_2.h5")
compile_model(model)

# Test
test_model(model)

# predict some of the pictures
img_names = os.listdir(S3_TO_PREDICT)
prefix = S3_TO_PREDICT

for img in img_names:
    print(img)
    img_tensor = load_image(prefix + img)
    pred = model.predict(img_tensor)[0][0]
    if pred >= 0.5:
        print("is a car")
    else:
        print("is not a car")
