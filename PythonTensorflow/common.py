import tensorflow as tf
from tensorflow.keras.layers import Conv2D, MaxPool2D, Flatten, Dense, Dropout
from keras.preprocessing import image
import numpy as np
import matplotlib.pyplot as plt
import os


LEARNING_RATE = 0.001  # start off with high rate first 0.001 and experiment with reducing it gradually

# The path to training data; may vary depending on actual machines
DATASET_PATH = 'D:/KaggleVehicleDetectionImageSet/train'
# The path to testing data; may vary depending on actual machines
TEST_DIR = 'D:/KaggleVehicleDetectionImageSet/test'

# The path where the images waiting to be predicted are cached
S3_TO_PREDICT = 'D:/Project2_August_2021/s3/toPredict'

DATA_LIST = os.listdir('D:/KaggleVehicleDetectionImageSet/train')

IMAGE_SIZE = (64, 64)  # size of the picture

NUM_CLASSES = len(DATA_LIST)
BATCH_SIZE = 40  # try reducing batch size or freeze more layers if your GPU runs out of memory
NUM_EPOCHS = 10

WEIGHTS_SAVE_PATH = "./weights/weight_3.h5"


def create_model():
    """
    Create and return a keras model.
    :return: a keras model.
    """
    conv_1 = Conv2D(input_shape=(64, 64, 3), filters=64, kernel_size=(3, 3), padding="same", activation="relu",
                    name="conv_1")
    conv_2 = Conv2D(filters=64, kernel_size=(3, 3), padding="same", activation="relu", name="conv_2")
    maxPool1 = MaxPool2D(pool_size=(2, 2), strides=(2, 2), name="maxPool1")

    conv_3 = Conv2D(filters=128, kernel_size=(3, 3), padding="same", activation="relu", name="conv_3")
    conv_4 = Conv2D(filters=128, kernel_size=(3, 3), padding="same", activation="relu", name="conv_4")
    maxPool2 = MaxPool2D(pool_size=(2, 2), strides=(2, 2), name="maxPool2")

    flatten_1 = Flatten(name="flatten_1")
    dense_1 = Dense(512, activation="relu", name="dense_1")
    dropout_1 = Dropout(0.2, name="dropout_1")
    dense_2 = Dense(256, activation="relu", name="dense_2")
    dropout_2 = Dropout(0.2, name="dropout_2")
    dense_3 = Dense(1, activation="sigmoid", name="dense_3")

    model = tf.keras.Sequential([conv_1, conv_2, maxPool1, conv_3, conv_4, maxPool2, flatten_1, dense_1, dropout_1,
                                 dense_2, dropout_2, dense_3])
    model.build(input_shape=(None, 64, 64, 3))
    model.summary()
    return model


def compile_model(model):
    """
    Compile the keras model.
    :param model: the keras model.
    :return: None.
    """
    lr_schedule = tf.keras.optimizers.schedules.ExponentialDecay(
        initial_learning_rate=LEARNING_RATE,
        decay_steps=10000,
        decay_rate=0.9)

    # compile
    model.compile(optimizer=tf.keras.optimizers.Adam(learning_rate=lr_schedule),
                  loss='binary_crossentropy',
                  metrics=['accuracy'])


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
