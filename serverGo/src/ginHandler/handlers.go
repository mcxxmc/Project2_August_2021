package ginHandler

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os"
	"serverGo/src/common"
	"serverGo/src/dbInterface"
	"serverGo/src/openCVGRPC"
	"serverGo/src/tensorflowGRPC"
	"time"
)

// The cache for labelImages; hold reference to the name and the path;
// will be useful when there are multiple labelers.
var mapNamesPaths = make(map[string]string)

// HandlerPostImage handles the uploaded image from the frontend.
func HandlerPostImage(c *gin.Context) {
	// IMPORTANT! allows CORS
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the file from the form using the key.
	file, _ := c.FormFile(common.FormFileName)

	if file == nil {
		c.JSON(http.StatusOK, gin.H{"r": "Please submit an image that is not empty."})
		return
	}

	// Checks if the image exists in our database.
	imgName := file.Filename
	exist, prediction, label, path := dbInterface.QueryName(imgName)

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
		dbInterface.InsertBared(imgName, imgPath)
	} else {
		// Else, checks the prediction and the label.
		msg = completeMsgIfNameExist(prediction, label, path)
	}

	// Make a response as a JSON object.
	c.JSON(http.StatusOK, gin.H{"r": msg})
}

// HandlerImmediatePred immediately predicts an image.
func HandlerImmediatePred(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// same codes in HandlerPostImage
	file, _ := c.FormFile(common.FormFileNameImmediatePred)
	if file == nil {
		c.JSON(http.StatusOK, gin.H{"r": "Please submit an image that is not empty."})
		return
	}
	imgName := file.Filename
	exist, prediction, label, path := dbInterface.QueryName(imgName)
	var msg string
	if exist == false {
		openedFile, err := file.Open()
		defer closeOpenedFile(openedFile)
		common.PanicErr(err)
		imgPath := common.S3ToPredict + imgName
		err = c.SaveUploadedFile(file, imgPath)
		common.PanicErr(err)
		dbInterface.InsertBared(imgName, imgPath)

		// gRPC client
		connectionTf, err := grpc.Dial(common.GRPCTensorflowPort, grpc.WithInsecure(), grpc.WithBlock())
		defer closeGRPCConnection(connectionTf)
		common.PanicErr(err)
		clientTf := tensorflowGRPC.NewCommunicatorClient(connectionTf)

		// send the request
		// TODO: maybe remove timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()
		r, err := clientTf.ImmediatePred(ctx, &tensorflowGRPC.Image{Name: imgName, Path: imgPath})
		common.PanicErr(err)

		predTf := r.Pred
		msg = completeMsgTensorflowResponse(predTf)
		newPath, err := getPathPredicted(imgName, predTf)
		if err == nil {
			err := os.Rename(imgPath, newPath)
			if err == nil {
				dbInterface.UpdatePathAndPrediction(imgName, newPath, predTf)
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

// HandlerOpenCV communicates with the OpenCV server.
func HandlerOpenCV(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// connects the OpenCV server through grpc
	connectionOpenCV, err := grpc.Dial(common.GRPCOpenCVInsecurePort, grpc.WithInsecure(), grpc.WithBlock())
	defer closeGRPCConnection(connectionOpenCV)
	common.PanicErr(err)
	clientOpenCV := openCVGRPC.NewCollectorClient(connectionOpenCV)

	// send the request
	// The timeout should be long enough; otherwise there will be errors
	// TODO: maybe remove timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	r, err := clientOpenCV.CollectImage(ctx, &openCVGRPC.Empty{})
	common.PanicErr(err)

	// update database
	dbInterface.InsertBared(r.Name, r.Path)

	c.Status(http.StatusOK)
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
				// If an image cannot be loaded, then skip it.
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
			mapNamesPaths[name] = path
		}else {
			// If an image cannot be loaded, then skip it.
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
							// TODO: maybe utilize the performance by using a single db connection
							// TODO: instead of creating one each time
							dbInterface.UpdatePathAndLabel(name, newLocation, isVehicle)
							// Remove the name from the map.
							delete(mapNamesPaths, name)
						}else {
							// Since we may have more than 1 user:
							fmt.Println(err)
						}
					} else {
						fmt.Println("LabelPicturesPost: Unexpected val.")
					}
				}
			}
		}
	}else {
		// Panic here to let the user know that there is a problem with the system.
		panic(err)
	}
	c.Status(http.StatusOK)
}
