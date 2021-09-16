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
	router.Use(webservice.Filter())
	router.Use(webservice.SetHeader())

	router.POST("/upload", webservice.PostImage)
	router.POST("/prediction", webservice.ImmediatePred)
	router.POST("/pictures", webservice.ShowPictures)
	router.POST("/labels-pictures", webservice.PostImageLabels)

	router.GET("/list", webservice.ShowList)
	router.GET("/labels-pictures", webservice.GetUnlabeledPictures)
	// TODO: safety / authorization
	router.GET("/opencv", webservice.UseCamera)

	// Run at port 8080.
	err := router.Run(common.GINPORT)
	common.CheckErr(err)
}
