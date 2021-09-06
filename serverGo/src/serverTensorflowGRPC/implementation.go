package serverTensorflowGRPC

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"serverGo/src/common"
	"serverGo/src/dbInterface"
	"serverGo/src/tensorflowGRPC"
)

type tfServer struct {
	tensorflowGRPC.UnimplementedCommunicatorServer
}

var mapNamesPaths = make(map[string]string)  // the dictionary storing old name-path pairs

func (s *tfServer) RequestImages(ctx context.Context, in *tensorflowGRPC.TFStandard) (*tensorflowGRPC.ImageArray, error) {
	fmt.Println("RequestImages: invoked.")
	records := dbInterface.FetchUnpredictedUnlabeled()
	var r tensorflowGRPC.ImageArray
	var tmp []*tensorflowGRPC.Image
	var name string
	var path string
	for i := 0; i < len(records) ; i ++ {
		record := records[i]
		name = record.Name
		path = record.Path
		tmp = append(tmp, &tensorflowGRPC.Image{Name: name, Path: path})
		mapNamesPaths[name] = path
	}
	r.Images = tmp
	return &r, nil
}

func (s *tfServer) PostPredictions(ctx context.Context, in *tensorflowGRPC.PredictionArray) (*tensorflowGRPC.TFStandard, error) {
	fmt.Println("PostPredictions: invoked.")
	predictions := in.Predictions
	var name string
	var b bool
	var path string
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
				dbInterface.UpdatePathAndPrediction(name, path, b)
				// remove the name from the map
				delete(mapNamesPaths, name)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println("PostPredictions: " + name + " does not exist in map.")
		}
	}
	var r tensorflowGRPC.TFStandard
	return &r, nil
}

func StartServer() {
	log.Printf("serverTensorflowGRPC: Server is started.")
	lis, err := net.Listen("tcp", common.GRPCGoPort)
	common.PanicErr(err)
	s := grpc.NewServer()
	tensorflowGRPC.RegisterCommunicatorServer(s, &tfServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
