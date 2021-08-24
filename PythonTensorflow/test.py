from tensorflow.keras.preprocessing.image import ImageDataGenerator
import numpy as np
from common import (
    TEST_DIR, IMAGE_SIZE
)


def test_model(model):
    """
    Test the model.
    :param model: a keras model.
    :return: None.
    """
    test_datagen = ImageDataGenerator(rescale=1. / 255)
    eval_generator = test_datagen.flow_from_directory(TEST_DIR, target_size=IMAGE_SIZE,
                                                      batch_size=1, shuffle=False, class_mode="binary")
    eval_generator.reset()
    x = model.evaluate_generator(eval_generator, steps=np.ceil(len(eval_generator)),
                                 use_multiprocessing=False, verbose=1, workers=1)
    print('Test loss: ', x[0])
    print('Test accuracy: ', x[1])
