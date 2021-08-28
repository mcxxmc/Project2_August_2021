from keras.preprocessing import image
import numpy as np
import matplotlib.pyplot as plt
from static.const import (
    S3_TO_PREDICT, S3_PRED_VEHICLE, S3_PRED_NON_VEHICLE
)
import os
import shutil


def load_image(img_path: str, show=False):
    """
    Load an image and return a numpy array.
    :param img_path: str.
    :param show: bool.
        Whether to display the image. Default is False.
    :return: a numpy array.
    """
    img = image.load_img(img_path, target_size=(64, 64))
    img_tensor = image.img_to_array(img)  # (height, width, channels)
    img_tensor = np.expand_dims(img_tensor, axis=0)  # (1, height, width, channels), add a dimension because the model
    # expects this shape: (batch_size, height, width, channels)
    img_tensor /= 255.  # imshow expects values in the range [0, 1]

    if show:
        plt.imshow(img_tensor[0])
        plt.axis('off')
        plt.show()

    return img_tensor


def make_prediction(model, folder_path: str = S3_TO_PREDICT, vehicle_path: str = S3_PRED_VEHICLE,
                    non_vehicle_path: str = S3_PRED_NON_VEHICLE):
    """
    Make predictions for images in the S3 toPredict folder, and then move them to the predicted folder accordingly.
    :param model: the keras model.
    :param folder_path: str.
        The path of the S3 toPredict folder.
    :param vehicle_path: str.
        The path of the S3 Predicted vehicle folder.
    :param non_vehicle_path: str.
        The path of the S3 Predicted non-vehicle folder.
    :return: .
    """
    img_names = os.listdir(folder_path)

    for img in img_names:
        img_path = folder_path + img
        print(img_path)
        img_tensor = load_image(img_path)
        pred = model.predict(img_tensor)[0][0]
        if pred >= 0.5:
            print("is a car")
            shutil.move(img_path, vehicle_path + img)
        else:
            print("is not a car")
            shutil.move(img_path, non_vehicle_path + img)

    # TODO: make a JSON (imgName: prediction, path) response to the Go server.
