package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"mime/multipart"
	"net/http"
	"projectGo/src/projectGo"
)

var channelsIO = 2  // number of threads for sql and for saving image to disk
var formFileName = "img"

var buffer = make(chan projectGo.BufferedImage)


// stores the image in the channel.
func bufferImage(c chan<- projectGo.BufferedImage, img projectGo.BufferedImage) {
	c <- img
	fmt.Println("A new image is in buffer.")
}

// writes the buffered image to disk
func flushBufferedImage(c <-chan projectGo.BufferedImage) {
	for {
		bfImage := <- c
		imgName := bfImage.ImgName
		projectGo.SaveImage(imgName, bfImage.Img) // save the image to disk

		//TODO: connect to tensorflow server to check the image

		projectGo.Insert(imgName, true) // insert the record into the database
		//TODO: move this to handler function when tensorflow server is on
	}
}

func handlerGet(c *gin.Context) {
	// return a html page
	c.HTML(http.StatusOK, "index.html", nil)
}

func handlerPostImage(c *gin.Context) {
	fmt.Println("HandlerPostImage invoked.")

	// get the file from the form using the name of the key
	file, err := c.FormFile(formFileName)
	projectGo.CheckErr(err)

	// check if the image exists in our database
	imgName := file.Filename
	exist, b, path := projectGo.QueryName(imgName)

	var msg string

	// if the image is never seen before
	if exist == false {
		msg = "The picture does not exist in the database."

		openedFile, err := file.Open()
		defer func(openedFile multipart.File) {
			err := openedFile.Close()
			projectGo.CheckErr(err)
		}(openedFile)
		projectGo.CheckErr(err)

		img, _, err := image.Decode(openedFile)
		projectGo.CheckErr(err)

		bufferImage(buffer, projectGo.BufferedImage{ImgName: imgName, Img: img})
		// the line above can be replaced by:
		// err = c.SaveUploadedFile(form.File, projectGo.PicTargetDir+imgName)
	} else {
		if b == true {
			msg = "true\n" + path
		}else {
			msg = "false\n" + path
		}
	}

	// make a response as a html file
	c.JSON(http.StatusOK, gin.H{
		"r": msg,
	})

	fmt.Println("HandlerPostImage finished.")
}

func main()  {
	projectGo.TryConnection()

	for i := 0; i < channelsIO; i++ {
		go flushBufferedImage(buffer)
	}

	router := gin.Default()
	router.LoadHTMLGlob("html/*.html")
	router.Static("/javascript", "./javascript")

	router.GET("/", handlerGet)
	router.POST("/imgSystem/", handlerPostImage)

	err := router.Run(":8080")  // run at port 8080
	projectGo.CheckErr(err)
}
