# Project2_August_2021
This is the second formal project I have completed out of school, starting august 2021. 

It is planned to be used as an add-on to my another CRUD project.

It is intended to build a system for image classification, which can be used for things like flaw detection on car components, object recognition, etc. The codes posted here are for illustration purpose and the examples used are for vehicle detection. The dataset for experiment is from Kaggle and the image size is 64 x 64 x 3.

## General

Keywords: Restful API, Micro service, gRPC, MVC

This project consists of frontend and backend. 

The frontend is implemented in html and javascript. Used JQuery and AJax.

The backend consists of 3 parts: opencv, tf-service, webserver. All 3 parts are independent and should be run as different servers, or at least as 5 different containers / images. (1 for opencv, 2 for tf-service and 2 for webserver, see details below) 

Also, there should be a MySQL server as the database.

* "opencv" is implemented in Python. It is used to invoke the camera and take a picture.

* "tf-service" is implemented in Python. It is used to predict the labels of pictures. It consists 2 parts: one is for predicting pictures periodically, the other is for immediate prediction. Each part has its own main file, so there should be 2 instances of this server.

* "webserver" is implemented in Golang. It serves as the middleman between frontend and backend. It uses Gin framework for router, and gorm to interact with the database. It also uses zap log for logging. There are 2 main files, one for restful API (to interact with frontend) and the other one for gRPC (to interact with other servers). Therefore, there should be 2 instances of this server. The 2 instances will not communicate with each other as they undertake different tasks that are not overlapping. However, they do interact with the same database and the same "S3" system.

"s3": here the S3 system is simplified. It is only a folder which contains images. In practical use, this system should be an independent server. It contains training data (labeled), predicted data (predicted but not labeled) and unpredicted data (wait to be predicted by our tensorflow model). 

WARNING: the Dockerfile in this project is outdated and no longer works.

### Frontend
Implemented in simple html and javascript. For practice use, these static file should be served by engines like Nginx. An example of Nginx configuration is included, but be sure to adjust it when necessary.

There are 5 pages (except index) corresponding to 5 different functionalities. All use AJax to get / send data from / to the backend. 

* "labelPictures.html": for labeling pictures manually. It can handle situations such as empty response and already labeled pictures (if there are more than 1 labeler at the same time). It will cause the pictures to be moved into different folders in S3.

* "openCV.html": for taking a picture using the camera. The picture will be put into S3.

* "showList.html": show the info of all the pictures in the database.

* "showPictures.html": show the info, including the visualizations themselves, of the picture. To avoid heavy load on internet, a fixed number of pictures (instead of all) are displayed on the screen.

* "uploadImg.html": for uploading image to S3. If the picture is already in the database, its info will pump up; otherwise the image will go to S3 and waits for prediction. In this project, image is only recognized by the name; however, this attribute can be easily changed. For real use, image hash value can be used to avoid this shortcoming. This page also consists of 2 parts. First is for normal upload, the other is for immediate prediction (which will predict the image immediately if its info is not in database).

Note: all changes made to the S3 will be reflected in the database immediately, so sometimes the prediction and label of a picture can both be "None" as it is just waiting to be predicted. However, even if a picture is not predicted, it will have a valid path (both on disk and recorded in the database; can be changed if it is predicted or labeled later).

### Opencv
It runs as a gRPC server. Every time it is called, a new photo is taken, assigned a random name and is transformed into the required size (64 x 64 x 3 here).

The main package used is `cv2`.

Runs on port 50051 by default.

### tf-service
It runs a tensorflow CNN model. It is has 2 parts and therefore 2 independent main files, one is a gRPC client (`tf`) and the other one is a gRPC server (`tf_fast`). The gRPC codes for client and server is seperated for good practice.

The client plays an active role: it contacts the webserver every 30 seconds (by default) asking for images to predict (more exactly, the paths of the images stored in S3). Then, it sends the result back to the webserver through another active contact.

The server is a passive one (this is the characteristic of gRPC): it receives image (or in fact, its path) from webserver, immediately predicts it, and sends the result back.

Note it does not make changes to S3; it only reads it. Changes (such as, moving a picture) are to be made by webserver, using the result / prediction of tf-service.

The 2 parts run as 2 independent instances simultaneously. Both runs on port 50052 by defualt. To avoid possible conflict, should be run on two different ports when in practical use.

### webserver
It serves as the middleman between frontend and backend. 

It uses Gin framework for router, and gorm to interact with the database. It also uses zap log for logging. 

There are 2 main files, one for restful API (to interact with frontend) and the other one for gRPC (to interact with other servers). Therefore, there should be 2 instances of this server. The 2 instances will not communicate with each other as they undertake different tasks that are not overlapping. However, they do interact with the same database and the same "S3" system.

By default, restful one runs on port 8080, gRPC one runs on port 50050.

Important:

#### db:

* uses gorm as a middleware between MySQL

* has a singleton Db object, which is the connection pool (or more exactly, a *gorm.Db object). Should be initiated by calling `db.openDb()` and closed by `db.closeDb()`.

#### opencv:

* runs as a gRPC client that actively contacts the opencv server (which is a gRPC server).

#### proto:

* all proto files. For convenience, I generate codes using these files in the same directory, then move them to the correct locations.

* in proto files, `Empty` and `TFStandard` are basically the same thing. However, gRPC itself does not allow same names. Another Reason why I make such duplication: the two are in different proto files and the copies of these files are kept by 2 seperate servers, so it is not likely that they can share the same `Empty`. Nevertheless, there can be a better way such as using `import` from a third party (e.g., google).

#### tf & tf_impelement:

* runs as a passive gRPC server waiting for requests from the tf-service server (which is a client). The impelement of the gRPC server is in `tf-implement`.

#### tf_fast:

* runs as an active gRPC client that sends request to the tf-service server (which is a gRPC server intead of a client for this time). For good practice, the codes & proto files of `tf` and `tf-fast` is seperated.

#### webservice:

* VERY IMPORTANT: BE SURE TO CHECK THE REAL CODES FOR DETAILS!

* the "controller". Used Gin framwork. Responsible for routing, binding url (in a restful api way), middlewares, making responses to users, etc. 

* the one who interacts with the MySQL database using gorm. It also makes changes to the S3 folder (e.g., moving one image from 1 folder to another). All changes of paths of images will be updated in the database.

* consisting of many handlers for GET and POST HTTP requests. They involve database operations, ajax, gRPC, disk operations, etc.

* the "label pictures" part takes additional steps to avoid potential problems like empty response from  user and multiple labelers (whose work may overlap, and a picture may no longer exist in the old path). Personally, I like these designs as they add more security to the program. 

* implemented delayed shutdown (by default is 5 seconds). Allows this server to finish the "last-minute" job.

Note: there is also a createtable.sql which is the original schema for the table in the database.

## TODO 

* encapsulate the usage of zap log in this project into a new "log" package, so there is no need to make a direct call to the "zap" package every time. (e.g., zap.S()).

* instead of calling `defer db.CloseSharedDb()` in `main()`, consider using a go routine in `main()` that is blocked by a channel listening to the ending signal. And the go routine will then proceed and close the database connection pool. In this case, the initialization of the connection pool can be moved to `init()`.

* also, authenification. This can be easily immigrated from my previous work (the CRUD one, which uses token in header).