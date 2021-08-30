package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"projectGo/src/projectGo"
)

var formFileName = "img"

// S3ToPredict The path where the images waiting to be predicted are cached
var S3ToPredict = "D:/Project2_August_2021/s3/toPredict/"


// handler the uploaded image from the front end
func handlerPostImage(c *gin.Context) {
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
}

// handle the response from the tensorflow server
func handlerPredictedImage(c *gin.Context) {
	// TODO
}

// handle the request to show all the picture info in a list
func handlerShowList(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  // IMPORTANT! allows CORS
	records := projectGo.FetchAll()
	c.JSON(http.StatusOK, records)
}

// handle the request to show pictures
func handlerShowPictures(c * gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var temp projectGo.JSONShowPictures
	var imageBundles projectGo.ImageBundles
	err := c.BindJSON(&temp)
	if err == nil {
		paths := projectGo.FetchPathsN(temp.Offset, temp.N)
		for i:= 0; i < len(paths); i ++ {
			path := paths[i]
			file, err:= ioutil.ReadFile(path)
			if err == nil {
				imageBundles.Images = append(imageBundles.Images,
					projectGo.ImageBundle{EncodedImage: base64.StdEncoding.EncodeToString(file), Text: path})
			}else{
				fmt.Println(err)
			}
		}
	}else{
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, imageBundles)
}

func main()  {
	projectGo.TryConnection()

	router := gin.Default()

	// here the router is only responsible for the POST request, as the GET request is handled by the front end
	router.POST("/imgSystem/", handlerPostImage)

	// This url is for requests from the tensorflow server
	router.POST("/fromTensorflow/", handlerPredictedImage)

	// for showPictures
	router.POST("/showPictures/", handlerShowPictures)

	// This url is for showList Request
	router.GET("/showList/", handlerShowList)

	err := router.Run(":8080")  // run at port 8080
	projectGo.CheckErr(err)
}
