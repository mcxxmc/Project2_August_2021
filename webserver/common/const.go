package common

// FormFileName The name of the form data in the form
var FormFileName = "img"

// FormFileNameImmediatePred The name of the form data in the form for immediate prediction.
var FormFileNameImmediatePred = "img_fast"

// S3ToPredict The path where the images waiting to be predicted are cached
var S3ToPredict = "D:/Project2_August_2021/s3/toPredict/"

// S3VehiclePredictionPrefix The path prefix for vehicle images,predicted and not labeled
var S3VehiclePredictionPrefix = "D:/Project2_August_2021/s3/predicted/vehicles/"

// S3NonVehiclePredictionPrefix The path prefix for non-vehicle images,predicted and not labeled
var S3NonVehiclePredictionPrefix = "D:/Project2_August_2021/s3/predicted/non-vehicles/"

// S3VehiclePrefix The path prefix for vehicle images,labeled
var S3VehiclePrefix = "D:/Project2_August_2021/s3/train/vehicles/"

// S3NonVehiclePrefix The path prefix for non-vehicle images,labeled
var S3NonVehiclePrefix = "D:/Project2_August_2021/s3/train/non-vehicles/"

// ResultIsVehicle string to match if the result is vehicle
var ResultIsVehicle = "v"

// ResultIsNonVehicle string to match if the result is non-vehicle
var ResultIsNonVehicle = "nv"

// WebserverPort the port number for webserver gRPC server
var WebserverPort = ":50050"

// OpenCVInsecurePort the port number for opencv gRPC server
var OpenCVInsecurePort = "localhost:50051"

// TensorflowPort the port number for tensorflow gRPC port
var TensorflowPort = "localhost:50052"

// WebserverPortGin the port number for webserver GIN server
var WebserverPortGin = ":8080"
