import cv2 as cv
import numpy
import matplotlib.pyplot as plt
import numpy as np

from const.const import (
    S3_TO_PREDICT, chars
)
import random


def capture(camera_seq: int = 0) -> numpy.ndarray:
    """
    Uses the device camera to take a picture and return as a numpy n-d array.
    :param camera_seq: int.
        The sequence number of the camera on the device. Default is 0. May vary depends on machine.
    :return: numpy.ndarray.
    """
    camera = cv.VideoCapture(camera_seq)  # access the webcam

    frame = None

    if camera.isOpened():
        _, frame = camera.read()

    camera.release()

    return frame


def bgr2rgb(frame: numpy.ndarray) -> numpy.ndarray:
    """
    Convert a frame from BGR to RGB.
    :param frame: numpy.ndarray.
    :return: numpy.ndarray.
    """
    return frame[:, :, ::-1]


def display(convertedFrame: numpy.ndarray) -> None:
    """
    Display the converted frame as an image.
    :param convertedFrame: numpy.ndarray.
    :return: None.
    """
    plt.imshow(convertedFrame)
    plt.show()


def generate_random_png_name() -> str:
    """
    Returns a randomly generated name for png file.
    :return: str.
        e.g., 'camera_abcdefgh123.png'
    """
    return 'camera_' + ''.join(random.sample(chars, 10)) + '.png'


def reshape_image_smaller(img: numpy.ndarray, new_height: int = 64, new_width: int = 64,
                          pooling: str = "mean") -> numpy.ndarray:
    """
    Reshape the image to make it smaller by pooling.
    :param img: numpy.ndarray.
    :param new_height: int. Default is 64.
    :param new_width: int. Default is 64.
    :param pooling: str.
        Default is "mean". Other methods include "min", "max", "median", "mid".
    :return: numpy.ndarray.
    """
    r = numpy.zeros((new_height, new_width, 3), dtype=int)

    old_height, old_width, rgb = img.shape

    step_size_height = old_height // new_height
    step_size_width = old_width // new_width

    f = False

    if pooling == "min":
        method = np.min
    elif pooling == "max":
        method = np.max
    elif pooling == "median":
        method = np.median
    elif pooling == "mid":
        method = None  # implements this special method in the nested for loop
        f = True
    else:
        method = np.mean

    for i in range(new_height):
        for j in range(new_width):
            if f:
                r[i][j] = img[i * step_size_height, j * step_size_width, ::]
            else:
                r[i][j] = method(img[i * step_size_height: (i + 1) * step_size_height,
                                  j * step_size_width: (j + 1) * step_size_width, ::])
    return r


def save_image(reshaped_img: numpy.ndarray, folder: str = S3_TO_PREDICT) -> (str, str):
    """
    Save the image to disk. Returns the image name and the path.
    :param reshaped_img: numpy.ndarray.
    :param folder: str. With default.
    :return: (str, str).
    """
    name = generate_random_png_name()
    imgPath = folder + name
    plt.imsave(imgPath, reshaped_img.astype(np.uint8))
    return name, imgPath
