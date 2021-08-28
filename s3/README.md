The S3 folder to hold pictures.
Can be deplyoed as an independent server.

Structure:

train/
    vehicles/
    non-vehicles/
toPredict/
predicted/
    vehicles/
    non-vehicles/

The "train" folder holds the pictures for training, which are classified into 2 classes.
The "toPredict" folder caches the pictures waiting to be classified by the tensorflow model. Once it is classified, it will be moved into "predicted" folder.
The "predicted" folder holds predicted pictures by the tensorflow model. Once those pictures are verified and labeled by human, they will be moved into "train" folder. 