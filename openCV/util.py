import cv2 as cv
import numpy
import matplotlib.pyplot as plt


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
