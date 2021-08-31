package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"projectGo/src/projectGo"
)

var formFileName = "img"

// S3ToPredict The path where the images waiting to be predicted are cached
var S3ToPredict = "D:/Project2_August_2021/s3/toPredict/"

// S3VehiclePrefix The path prefix for vehicle images,labeled
var S3VehiclePrefix = "D:/Project2_August_2021/s3/train/vehicles/"

// S3NonVehiclePrefix The path prefix for non-vehicle images,labeled
var S3NonVehiclePrefix = "D:/Project2_August_2021/s3/train/non-vehicles/"

var resultVehicle = "v"

// The cache for labelImages; hold reference to the name and the path; will be useful when there are
// multiple labelers
var name2path = make(map[string]string)


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

		err = c.SaveUploadedFile(file, S3ToPredict + imgName)
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
		pathsDescs := projectGo.FetchN(temp.Offset, temp.N)
		var path string
		var text string
		for i := 0; i < len(pathsDescs); i ++ {
			path = pathsDescs[i].Path
			text = pathsDescs[i].Text
			file, err := ioutil.ReadFile(path)
			if err == nil {
				imageBundles.Images = append(imageBundles.Images,
					projectGo.ImageBundle{EncodedImage: base64.StdEncoding.EncodeToString(file), Text: text})
			}else{
				fmt.Println(err)
			}
		}
	}else{
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, imageBundles)
}

// handle request to label the pictures; GET
func handlerLabelPicturesGET(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var imageBundles projectGo.ImageBundles
	unlabeledRecords := projectGo.FetchUnlabeled()
	var name string
	var path string
	for i := 0; i < len(unlabeledRecords); i ++ {
		name = unlabeledRecords[i].Name
		path = unlabeledRecords[i].Path
		file, err := ioutil.ReadFile(path)
		if err == nil {
			imageBundles.Images = append(imageBundles.Images,
				projectGo.ImageBundle{EncodedImage: base64.StdEncoding.EncodeToString(file), Text: name})
			// cache the name-path pair
			name2path[name] = path
		}else {
			fmt.Println(err)
		}
	}
	c.JSON(http.StatusOK, imageBundles)
}

// handle request to label the pictures; POST
func handlerLabelPicturesPOST(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var temp projectGo.JSONLabeledResults
	err := c.BindJSON(&temp)
	if err == nil {
		results := temp.Results
		// first check if the results contain information
		if results == nil || len(results) == 0 {
			fmt.Println("LabelPicturesPost: Empty response from users.")
		}else{
			var newLocation string
			var isVehicle bool
			for i := 0; i < len(results); i ++ {
				result := results[i]
				name := result.Name
				path, exist := name2path[name]
				// ignore if there is no such name; maybe the record has been updated by another user
				// else, move the pictures to the correct folder.
				if exist == true {
					val := result.Val
					if val == resultVehicle {
						newLocation = S3VehiclePrefix + name
						isVehicle = true
					}else {
						newLocation = S3NonVehiclePrefix + name
						isVehicle = false
					}
					err := os.Rename(path, newLocation)
					if err == nil {
						// also, update the database
						projectGo.UpdatePathAndLabel(name, newLocation, isVehicle)
						// remove the name from the map
						delete(name2path, name)
					}else {
						fmt.Println(err)
					}
				}
			}
		}
	}else {
		fmt.Println(err)
	}
	c.Status(http.StatusOK)
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

	// for labelPictures
	router.GET("/labelPictures/", handlerLabelPicturesGET)
	router.POST("/labelPictures/", handlerLabelPicturesPOST)

	err := router.Run(":8080")  // run at port 8080
	projectGo.CheckErr(err)
}
