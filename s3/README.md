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


For convenience (should be removed later):
In MySQL workBench:
INSERT INTO picture(name, path, prediction, label) values ("1.png", "D:/Project2_August_2021/s3/train/vehicles/1.png", null, true);
INSERT INTO picture(name, path, prediction, label) values ("2.png", "D:/Project2_August_2021/s3/train/vehicles/2.png", null, true);
INSERT INTO picture(name, path, prediction, label) values ("3.png", "D:/Project2_August_2021/s3/train/vehicles/3.png", null, true);
INSERT INTO picture(name, path, prediction, label) values ("4.png", "D:/Project2_August_2021/s3/train/vehicles/4.png", null, true);
INSERT INTO picture(name, path, prediction, label) values ("5.png", "D:/Project2_August_2021/s3/train/vehicles/5.png", null, true);
INSERT INTO picture(name, path, prediction, label) values ("extra1.png", "D:/Project2_August_2021/s3/train/non-vehicles/extra1.png", null, false);
INSERT INTO picture(name, path, prediction, label) values ("extra2.png", "D:/Project2_August_2021/s3/train/non-vehicles/extra2.png", null, false);
INSERT INTO picture(name, path, prediction, label) values ("extra3.png", "D:/Project2_August_2021/s3/train/non-vehicles/extra3.png", null, false);
INSERT INTO picture(name, path, prediction, label) values ("extra4.png", "D:/Project2_August_2021/s3/train/non-vehicles/extra4.png", null, false);
INSERT INTO picture(name, path, prediction, label) values ("extra5.png", "D:/Project2_August_2021/s3/train/non-vehicles/extra5.png", null, false);