import tensorflow as tf
from tensorflow.keras.layers import Conv2D, MaxPool2D, Flatten, Dense, Dropout
from static.const import LEARNING_RATE


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
    # model.summary()
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
