from tensorflow.keras.preprocessing.image import ImageDataGenerator
import matplotlib.pyplot as plt

from model.model import create_model, compile_model
from static.const import (
    DATASET_PATH, IMAGE_SIZE, NUM_EPOCHS, BATCH_SIZE, WEIGHTS_SAVE_PATH
)


def train(model=None, save_weights: bool = True):
    """
    Train a model and return it.
    Optionally save the weights.
    Configurations including dataset paths and image size are set in and import from model.py.
    :param model: a keras model.
        Default is None. If None, a new model will be created.
        If a model is to be passed in using this parameter, the model passed in should not have been compiled.
    :param save_weights: bool.
        Default is True. Whether to save the weights.
    :return: a keras model.
    """
    # Generating Training and Validation Batches
    # with noise
    train_datagen = ImageDataGenerator(rescale=1./255, rotation_range=50, featurewise_center=True,
                                       featurewise_std_normalization=True, width_shift_range=0.2,
                                       height_shift_range=0.2, shear_range=0.25, zoom_range=0.1,
                                       zca_whitening=True, channel_shift_range=20,
                                       horizontal_flip=True, vertical_flip=True,
                                       validation_split=0.2, fill_mode='constant')

    train_batches = train_datagen.flow_from_directory(DATASET_PATH, target_size=IMAGE_SIZE,
                                                      shuffle=True, batch_size=BATCH_SIZE,
                                                      subset="training", seed=42,
                                                      class_mode="binary")

    valid_batches = train_datagen.flow_from_directory(DATASET_PATH, target_size=IMAGE_SIZE,
                                                      shuffle=True, batch_size=BATCH_SIZE,
                                                      subset="validation", seed=42,
                                                      class_mode="binary")
    if model is None:
        # build the model
        model = create_model()

    compile_model(model)

    # Train
    print("number of train batches is: " + str(len(train_batches)))
    print("number of valid batches is: " + str(len(valid_batches)))

    STEP_SIZE_TRAIN = train_batches.n//train_batches.batch_size
    STEP_SIZE_VALID = valid_batches.n//valid_batches.batch_size

    print("Step size train is: " + str(STEP_SIZE_TRAIN))
    print("Step size validation is: " + str(STEP_SIZE_VALID))

    history = model.fit_generator(generator=train_batches, epochs=NUM_EPOCHS, validation_data=valid_batches,
                                  steps_per_epoch=STEP_SIZE_TRAIN, validation_steps=STEP_SIZE_VALID)

    if save_weights is True:
        # Save weights
        model.save_weights(WEIGHTS_SAVE_PATH)

    # Plot the loss graph
    plt.plot(history.history['accuracy'], label='Train_acc')
    plt.plot(history.history['val_accuracy'], label='val_acc')
    plt.xlabel('Accuracy over 10 Epochs')
    plt.legend(loc='lower right')
    plt.grid(True)
    plt.show()

    plt.plot(history.history['loss'], label='Train_loss')
    plt.plot(history.history['val_loss'], label='val_loss')
    plt.xlabel('Loss over 10 Epochs')
    plt.legend(loc='upper right')
    plt.grid(True)
    plt.show()

    return model
