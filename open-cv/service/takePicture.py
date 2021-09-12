from util.functions import (
    capture, bgr2rgb, display, reshape_image_smaller, save_image
)
from const.const import (
    CAMERA_SEQ, DEFAULT_POOLING
)


def take_picture() -> (str, str):
    """
    The function call to use the function to take a picture and save it in S3.
    Returns the image name and its path.
    :return: (str, str).
    """
    frame = capture(CAMERA_SEQ)
    frameRGB = bgr2rgb(frame)  # BGR => RGB
    display(frameRGB)
    reshaped = reshape_image_smaller(frameRGB, pooling=DEFAULT_POOLING)
    display(reshaped)
    return save_image(reshaped)
