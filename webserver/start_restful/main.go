package main

import (
	"github.com/gin-gonic/gin"
	"webserver/common"
	"webserver/db"
	"webserver/webservice"
)

func main() {
	// Test the connection to the database. Create the table if none is available.
	db.TryConnection()

	// Set up the router.
	router := gin.Default()

	router.POST("/upload-image", webservice.HandlerPostImage)
	router.POST("/fast-prediction", webservice.HandlerImmediatePred)
	router.POST("/show-pictures", webservice.HandlerShowPictures)
	router.POST("/label-pictures", webservice.HandlerLabelPicturesPOST)

	router.GET("/show-list", webservice.HandlerShowList)
	router.GET("/label-pictures", webservice.HandlerLabelPicturesGET)
	// TODO: safety / authorization
	router.GET("/opencv", webservice.HandlerOpenCV)

	// Run at port 8080.
	err := router.Run(common.GINPORT)
	common.CheckErr(err)
}
