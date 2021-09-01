package ginHandler

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"serverGo/src/common"
	"serverGo/src/dbInterface"
)

// The cache for labelImages; hold reference to the name and the path;
// will be useful when there are multiple labelers.
var name2path = make(map[string]string)

// HandlerPostImage handles the uploaded image from the frontend.
func HandlerPostImage(c *gin.Context) {
	// IMPORTANT! allows CORS
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the file from the form using the key.
	file, err := c.FormFile(common.FormFileName)
	common.CheckErr(err)

	// Checks if the image exists in our database.
	imgName := file.Filename
	exist, prediction, label, path := dbInterface.QueryName(imgName)

	msg := ""

	// If the image is never seen before:
	if exist == false {
		msg = "The picture does not exist in the database."

		openedFile, err := file.Open()
		defer func(openedFile multipart.File) {
			err := openedFile.Close()
			common.CheckErr(err)
		}(openedFile)
		common.CheckErr(err)

		err = c.SaveUploadedFile(file, common.S3ToPredict + imgName)
	} else {
		// Else, checks the prediction and the label.
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

	// Make a response as a JSON object.
	c.JSON(http.StatusOK, gin.H{
		"r": msg,
	})
}

// HandlerPredictedImage handles the request from the tensorflow server.
func HandlerPredictedImage(c *gin.Context) {
	// TODO
}

// HandlerOpenCV handles the request from the OpenCV server.
func HandlerOpenCV(c *gin.Context) {
	// TODO
}

// HandlerShowList handles the request to show information of all the images in a list.
func HandlerShowList(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	records := dbInterface.FetchAll()
	c.JSON(http.StatusOK, records)
}

// HandlerShowPictures handles the request to show pictures.
func HandlerShowPictures(c * gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var temp JSONShowPictures
	var imageBundles ImageBundles
	err := c.BindJSON(&temp)
	if err == nil {
		pathsDescs := dbInterface.FetchN(temp.Offset, temp.N)
		var path string
		var text string
		for i := 0; i < len(pathsDescs); i ++ {
			path = pathsDescs[i].Path
			text = pathsDescs[i].Text
			file, err := ioutil.ReadFile(path)
			if err == nil {
				imageBundles.Images = append(imageBundles.Images,
					ImageBundle{EncodedImage: base64.StdEncoding.EncodeToString(file), Text: text})
			}else{
				fmt.Println(err)
			}
		}
	}else{
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, imageBundles)
}

// HandlerLabelPicturesGET handles the request to label the pictures; GET.
func HandlerLabelPicturesGET(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var imageBundles ImageBundles
	unlabeledRecords := dbInterface.FetchUnlabeled()
	var name string
	var path string
	for i := 0; i < len(unlabeledRecords); i ++ {
		name = unlabeledRecords[i].Name
		path = unlabeledRecords[i].Path
		file, err := ioutil.ReadFile(path)
		if err == nil {
			imageBundles.Images = append(imageBundles.Images,
				ImageBundle{EncodedImage: base64.StdEncoding.EncodeToString(file), Text: name})

			// Caches the name-path pair.
			name2path[name] = path
		}else {
			fmt.Println(err)
		}
	}
	c.JSON(http.StatusOK, imageBundles)
}

// HandlerLabelPicturesPOST handles the request to label the pictures; POST.
func HandlerLabelPicturesPOST(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var temp JSONLabeledResults
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
			for i := 0; i < len(results); i ++ {
				result := results[i]
				name := result.Name
				path, exist := name2path[name]
				// Ignore if there is no such name; maybe the record has been updated by another user.
				// Else, move the pictures to the correct folder.
				// Be careful if the val is empty.
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
							dbInterface.UpdatePathAndLabel(name, newLocation, isVehicle)
							// Remove the name from the map.
							delete(name2path, name)
						}else {
							fmt.Println(err)
						}
					} else {
						fmt.Println("LabelPicturesPost: Unexpected val.")
					}
				}
			}
		}
	}else {
		fmt.Println(err)
	}
	c.Status(http.StatusOK)
}
