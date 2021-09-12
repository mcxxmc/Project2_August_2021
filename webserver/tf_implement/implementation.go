package tf_implement

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"webserver/common"
	"webserver/db"
	tf2 "webserver/tf"
)

type tfServer struct {
	tf2.UnimplementedCommunicatorServer
}

var mapNamesPaths = make(map[string]string)  // the dictionary storing old name-path pairs

func (s *tfServer) RequestImages(ctx context.Context, in *tf2.TFStandard) (*tf2.ImageArray, error) {
	fmt.Println("RequestImages: invoked.")
	dbc := db.OpenDb()
	defer db.CloseDb(dbc)
	records := db.FetchUnpredictedUnlabeled(dbc).Recs
	var r tf2.ImageArray
	var tmp []*tf2.Image
	var name string
	var path string
	for i := 0; i < len(records) ; i ++ {
		record := records[i]
		name = record.Name
		path = record.Path
		tmp = append(tmp, &tf2.Image{Name: name, Path: path})
		mapNamesPaths[name] = path
	}
	r.Images = tmp
	return &r, nil
}

func (s *tfServer) PostPredictions(ctx context.Context, in *tf2.PredictionArray) (*tf2.TFStandard, error) {
	fmt.Println("PostPredictions: invoked.")
	predictions := in.Predictions
	var name string
	var b bool
	var path string
	dbc := db.OpenDb()
	defer db.CloseDb(dbc)
	// use a loop to process each prediction
	for i := 0; i < len(predictions); i ++ {
		pred := predictions[i]
		name = pred.Name
		oldPath, exist := mapNamesPaths[name]
		// check if the image name is still valid
		if exist == true {
			b = pred.Pred
			if b == true {
				path = common.S3VehiclePredictionPrefix + name
			} else if b == false {
				path = common.S3NonVehiclePredictionPrefix + name
			} else {
				fmt.Println("PostPredictions: Unexpected value b.")
				break
			}
			// fmt.Println("Name: " + name + ", old path: " + oldPath + ", path: " + path + ", pred: " + strconv.FormatBool(b))
			// move the image
			err := os.Rename(oldPath, path)
			if err == nil {
				// updates the database
				db.UpdatePathAndPrediction(dbc, name, path, b)
				// remove the name from the map
				delete(mapNamesPaths, name)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println("PostPredictions: " + name + " does not exist in map.")
		}
	}
	var r tf2.TFStandard
	return &r, nil
}

func StartServer() {
	log.Printf("tf_implement: Server is started.")
	lis, err := net.Listen("tcp", common.GRPCGoPort)
	common.PanicErr(err)
	s := grpc.NewServer()
	tf2.RegisterCommunicatorServer(s, &tfServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
