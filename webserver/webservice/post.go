package webservice

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"webserver/common"
	"webserver/db"
	tf2 "webserver/tf"
	tf2Fast "webserver/tf_fast"
)

// The cache for labelImages; hold reference to the name and the path;
// will be useful when there are multiple labelers.
var mapNamesPaths = make(map[string]string)

// PostImage handles the uploaded image from the frontend.
func PostImage(c *gin.Context) {
	// Get the file from the form using the key.
	file, _ := c.FormFile(common.FormFileName)

	if file == nil {
		c.JSON(http.StatusOK, gin.H{"r": "Please submit an image that is not empty."})
		return
	}

	dbc := db.OpenDb()
	defer db.CloseDb(dbc)

	// Checks if the image exists in our database.
	imgName := file.Filename
	exist, prediction, label, path := db.QueryName(dbc, imgName)

	var msg string

	// If the image is never seen before:
	if exist == false {
		msg = "The picture does not exist in the database.\n"

		openedFile, err := file.Open()
		defer closeOpenedFile(openedFile)
		common.PanicErr(err)

		imgPath := common.S3ToPredict + imgName
		err = c.SaveUploadedFile(file, imgPath)
		common.PanicErr(err)

		// update the database
		db.InsertBared(dbc, imgName, imgPath)

	} else {
		// Else, checks the prediction and the label.
		msg = completeMsgIfNameExist(prediction, label, path)
	}

	// Make a response as a JSON object.
	c.JSON(http.StatusOK, gin.H{"r": msg})
}

// ImmediatePred immediately predicts an image.
func ImmediatePred(c *gin.Context) {
	// same codes in PostImage
	file, _ := c.FormFile(common.FormFileNameImmediatePred)
	if file == nil {
		c.JSON(http.StatusOK, gin.H{"r": "Please submit an image that is not empty."})
		return
	}
	dbc := db.OpenDb()
	defer db.CloseDb(dbc)
	imgName := file.Filename
	exist, prediction, label, path := db.QueryName(dbc, imgName)
	var msg string
	if exist == false {
		openedFile, err := file.Open()
		defer closeOpenedFile(openedFile)
		common.PanicErr(err)
		imgPath := common.S3ToPredict + imgName
		err = c.SaveUploadedFile(file, imgPath)
		common.PanicErr(err)
		db.InsertBared(dbc, imgName, imgPath)

		// gRPC client
		connectionTf, err := grpc.Dial(common.GRPCTensorflowPort, grpc.WithInsecure(), grpc.WithBlock())
		defer closeGRPCConnection(connectionTf)
		common.PanicErr(err)
		clientTfFast := tf2Fast.NewImmediatePredictorClient(connectionTf)

		// send the request
		// TODO: maybe remove timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
		defer cancel()
		r, err := clientTfFast.ImmediatePred(ctx, &tf2.Image{Name: imgName, Path: imgPath})
		common.PanicErr(err)

		predTf := r.Pred
		msg = completeMsgTensorflowResponse(predTf)
		newPath, err := getPathPredicted(imgName, predTf)
		if err == nil {
			err := os.Rename(imgPath, newPath)
			if err == nil {
				db.UpdatePathAndPrediction(dbc, imgName, newPath, predTf)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		msg = completeMsgIfNameExist(prediction, label, path)
	}
	c.JSON(http.StatusOK, gin.H{"r": msg})
}

// ShowPictures handles the request to show pictures.
func ShowPictures(c * gin.Context) {
	var temp QueryParameters
	var imageBundles ImageBundles
	err := c.BindJSON(&temp)
	if err == nil {
		dbc := db.OpenDb()
		defer db.CloseDb(dbc)
		records := db.FetchN(dbc, temp.Offset, temp.N).Recs
		var path string
		var text string
		var record db.Record
		for i := 0; i < len(records); i ++ {
			record = records[i]
			path = record.Path
			text = generateText(record.Id, record.Name, record.Prediction, record.Label)
			file, err := ioutil.ReadFile(path)
			if err == nil {
				imageBundles.Images = append(imageBundles.Images,
					ImageBundle{EncodedImage: base64.StdEncoding.EncodeToString(file), Text: text})
			}else{
				// If an image cannot be loaded, then skip it.
				fmt.Println(err)
			}
		}
	}else{
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, imageBundles)
}

// PostImageLabels handles the request to label the pictures; POST.
func PostImageLabels(c *gin.Context){
	var temp LabeledResults
	err := c.BindJSON(&temp)
	if err == nil {
		results := temp.Results
		// First, check if the results contain information.
		if results == nil || len(results) == 0 {
			fmt.Println("LabelPicturesPost: Empty response from users.")
		}else{
			var newLocation string
			var isVehicle bool
			var hasVal = false
			dbc := db.OpenDb()
			defer db.CloseDb(dbc)
			for i := 0; i < len(results); i ++ {
				result := results[i]
				name := result.Name
				path, exist := mapNamesPaths[name]
				// Ignore if there is no such name; maybe the record has been updated by another user.
				// Else, move the pictures to the correct folder. Be careful if the val is empty.
				if exist == true {
					val := result.Val
					if val == common.ResultIsVehicle {
						newLocation = common.S3VehiclePrefix + name
						isVehicle = true
						hasVal = true
					}else if val == common.ResultIsNonVehicle {
						newLocation = common.S3NonVehiclePrefix + name
						isVehicle = false
						hasVal = true
					}
					if hasVal == true {
						// Move the image to the corresponding folder.
						err := os.Rename(path, newLocation)
						if err == nil {
							// Also, update the database.
							db.UpdatePathAndLabel(dbc, name, newLocation, isVehicle)
							// Remove the name from the map.
							delete(mapNamesPaths, name)
						}else {
							// Since we may have more than 1 user:
							fmt.Println(err)
						}
					} else {
						fmt.Println("LabelPicturesPost: Unexpected val.")
					}
				} else {
					fmt.Println("The image name does not exist in the cache.")
				}
			}
		}
	}else {
		// Panic here to let the user know that there is a problem with the system.
		panic(err)
	}
	c.Status(http.StatusOK)
}
