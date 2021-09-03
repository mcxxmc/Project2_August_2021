package main

import (
	"github.com/gin-gonic/gin"
	"serverGo/src/common"
	"serverGo/src/dbInterface"
	"serverGo/src/ginHandler"
)


func main()  {
	// Test the connection to the database. Create the table if none is available.
	dbInterface.TryConnection()

	// Set up the router.
	router := gin.Default()

	// For Uploaded images. It will return the label of the image or put it in the "toPredict" folder
	// if the image has never been seen before.
	router.POST("/imgSystem/", ginHandler.HandlerPostImage)

	// For interaction with the tensorflow server. Receive predictions and encode them into the database.
	// TODO
	router.POST("/fromTensorflow/", ginHandler.HandlerPredictedImage)

	// For "OpenCV". When a user visits this link remotely, the camera will be used.
	// TODO: safety / authorization
	router.GET("/fromOpenCV/", ginHandler.HandlerOpenCV)

	// For "showList". Show all records in our database.
	router.GET("/showList/", ginHandler.HandlerShowList)

	// For "showPictures". Show the pictures and the corresponding description.
	router.POST("/showPictures/", ginHandler.HandlerShowPictures)

	// For "labelPictures". Enable the users to manually label the predicted images and reflect the result
	// in the database accordingly.
	router.GET("/labelPictures/", ginHandler.HandlerLabelPicturesGET)
	router.POST("/labelPictures/", ginHandler.HandlerLabelPicturesPOST)

	// Run at port 8080.
	err := router.Run(":8080")
	common.CheckErr(err)
}
