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


def make_prediction(model, imgPath: str) -> bool:
    """
    Make predictions for a given image and returns a bool.
    :param model: the keras model.
    :param imgPath: str.
        The path of the image to predict.
    :return: bool.
    """
    img_tensor = load_image(imgPath)
    pred = model.predict(img_tensor)[0][0]
    if pred >= 0.5:
        return True
    else:
        return False
