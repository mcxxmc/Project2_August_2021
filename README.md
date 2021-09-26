# Project2_August_2021
This is the second formal project I have completed out of school, august 2021. 

It is intended to build a system for image classification, which can be used for things like flaw detection on car components, object recognition, etc. The codes posted here are for illustration purpose and the examples used are for vehicle detection. The dataset for experiment is from Kaggle and the image size is 64 x 64 x 3.

## General

This project consists of frontend and backend. 

The frontend is implemented in html and javascript. Used JQuery and AJax.

The backend consists of 3 parts: opencv, tf-service, webserver. All 3 parts are independent and should be run as different servers, or at least as 5 different containers / images. (1 for opencv, 2 for tf-service and 2 for webserver, see details below) 

Also, there should be a MySQL server as the database.

* "opencv" is implemented in Python. It is used to invoke the camera and take a picture.

* "tf-service" is implemented in Python. It is used to predict the labels of pictures. It consists 2 parts: one is for predicting pictures periodically, the other is for immediate prediction. Each part has its own main file, so there should be 2 instances of this server.

* "webserver" serves as the middleman between frontend and backend. It uses Gin framework for router, and gnorm to interact with the database. It also uses zap log for logging. There are 2 main files, one for restful API (to interact with frontend) and the other one for gRPC (to interact with other servers). Therefore, there should be 2 instances of this server. The 2 instances will not communicate with each other as they undertake different tasks that are not overlapping. However, they do interact with the same database and the same "S3" system.

"s3": here the S3 system is simplified. It is only a folder which contains images. In practical use, this system should be an independent server. It contains training data (labeled), predicted data (predicted but not labeled) and unpredicted data (wait to be predicted by our tensorflow model). 





## TODO 

* encapsulate the usage of zap log in this project into a new "log" package, so there is no need to make a direct call to the "zap" package every time. (e.g., zap.S()).

* instead of calling `defer db.CloseSharedDb()` in `main()`, consider using a go routine in `main()` that is blocked by a channel listening to the ending signal. And the go routine will then proceed and close the database connection pool. In this case, the initialization of the connection pool can be moved to `init()`.