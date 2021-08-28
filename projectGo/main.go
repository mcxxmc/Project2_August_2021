package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"projectGo/src/projectGo"
)

var formFileName = "img"

// S3ToPredict The path where the images waiting to be predicted are cached
var S3ToPredict = "D:/Project2_August_2021/s3/toPredict/"


// handler the uploaded image from the front end
func handlerPostImage(c *gin.Context) {
	fmt.Println("HandlerPostImage invoked.")

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  // IMPORTANT! allows CORS

	// get the file from the form using the name of the key
	file, err := c.FormFile(formFileName)
	projectGo.CheckErr(err)

	// check if the image exists in our database
	imgName := file.Filename
	exist, prediction, label, path := projectGo.QueryName(imgName)

	msg := ""

	// if the image is never seen before
	if exist == false {
		msg = "The picture does not exist in the database."

		openedFile, err := file.Open()
		defer func(openedFile multipart.File) {
			err := openedFile.Close()
			projectGo.CheckErr(err)
		}(openedFile)
		projectGo.CheckErr(err)

		err = c.SaveUploadedFile(file, S3ToPredict+ imgName)
	} else {  //else, check the prediction and the label
		if prediction != nil {
			if *prediction == true {
				msg += "prediction true\n"
			}else {
				msg += "prediction false\n"
			}
		}

		if label != nil {
			if *label == true {
				msg += "label true\n"
			}else {
				msg += "label false\n"
			}
		}

		msg += *path
	}

	// make a response as a html file
	c.JSON(http.StatusOK, gin.H{
		"r": msg,
	})

	fmt.Println("HandlerPostImage finished.")
}

// handle the response from the tensorflow server
func handlerPredictedImage(c *gin.Context) {}

func main()  {
	projectGo.TryConnection()

	router := gin.Default()

	// here the router is only responsible for the POST request, as the GET request is handled by the front end
	router.POST("/imgSystem/", handlerPostImage)

	// This url is for requests from the tensorflow server
	router.POST("/fromTensorflow/", handlerPredictedImage)

	err := router.Run(":8080")  // run at port 8080
	projectGo.CheckErr(err)
}
