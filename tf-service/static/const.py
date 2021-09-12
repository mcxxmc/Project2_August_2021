import os

LEARNING_RATE = 0.001  # start off with high rate first 0.001 and experiment with reducing it gradually

# The path to training data; may vary depending on actual machines
DATASET_PATH = 'D:/KaggleVehicleDetectionImageSet/train'
# The path to testing data; may vary depending on actual machines
TEST_DIR = 'D:/KaggleVehicleDetectionImageSet/test'

# The path where the images waiting to be predicted are cached
S3_TO_PREDICT = 'D:/Project2_August_2021/s3/toPredict/'

# The path where the predicted vehicle images go
S3_PRED_VEHICLE = 'D:/Project2_August_2021/s3/predicted/vehicles/'

# The path where the predicted non-vehicle images go
S3_PRED_NON_VEHICLE = 'D:/Project2_August_2021/s3/predicted/non-vehicles/'

DATA_LIST = os.listdir(DATASET_PATH)

IMAGE_SIZE = (64, 64)  # size of the picture

NUM_CLASSES = len(DATA_LIST)
BATCH_SIZE = 40  # try reducing batch size or freeze more layers if your GPU runs out of memory
NUM_EPOCHS = 10

WEIGHTS_PATH = '../weights/weight_2.h5'
WEIGHTS_SAVE_PATH = "../weights/weight_3.h5"

GRPC_INSECURE_PORT = "[::]:50052"
GRPC_GOLANG_INSECURE_PORT = "localhost:50050"
