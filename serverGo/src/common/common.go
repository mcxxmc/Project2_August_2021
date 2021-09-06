package common

import "fmt"

// FormFileName The name of the form data in the form
var FormFileName = "img"

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

var ResultIsVehicle = "v"
var ResultIsNonVehicle = "nv"

var GRPCGoPort = ":50050"
var GRPCOpenCVInsecurePort = "localhost:50051"
var GRPCTensorflowPort = "localhost:50052"

// CheckErr checks the error and prints it out if not nil
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// PanicErr raises a panic if the error is not nil.
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
