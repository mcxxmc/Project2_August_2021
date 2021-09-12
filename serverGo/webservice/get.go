package webservice

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"serverGo/common"
	"serverGo/db"
	opencv2 "serverGo/opencv"
	"time"
)

// HandlerOpenCV communicates with the OpenCV server.
func HandlerOpenCV(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// connects the OpenCV server through grpc
	connectionOpenCV, err := grpc.Dial(common.GRPCOpenCVInsecurePort, grpc.WithInsecure(), grpc.WithBlock())
	defer closeGRPCConnection(connectionOpenCV)
	common.PanicErr(err)
	clientOpenCV := opencv2.NewCollectorClient(connectionOpenCV)

	// send the request
	// The timeout should be long enough; otherwise there will be errors
	// TODO: maybe remove timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	r, err := clientOpenCV.CollectImage(ctx, &opencv2.Empty{})
	common.PanicErr(err)

	// update database
	dbc := db.OpenDb()
	defer db.CloseDb(dbc)
	db.InsertBared(dbc, r.Name, r.Path)

	c.Status(http.StatusOK)
}

// HandlerShowList handles the request to show information of all the images in a list.
func HandlerShowList(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	dbc := db.OpenDb()
	defer db.CloseDb(dbc)
	records := db.FetchAll(dbc)
	c.JSON(http.StatusOK, records)
}

// HandlerLabelPicturesGET handles the request to label the pictures; GET.
func HandlerLabelPicturesGET(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var imageBundles ImageBundles
	dbc := db.OpenDb()
	defer db.CloseDb(dbc)
	unlabeledRecords := db.FetchUnlabeled(dbc).Recs
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
