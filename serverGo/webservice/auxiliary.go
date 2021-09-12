package webservice

import (
	"errors"
	"google.golang.org/grpc"
	"mime/multipart"
	"serverGo/common"
	"strconv"
)

// ImageBundle The JSON structure to store an encoded base64 image and the corresponding text.
// Content of the text may vary depending on different usages.
type ImageBundle struct {
	EncodedImage string `json:"image"`
	Text string `json:"text"`
}

// ImageBundles A collection of ImageBundle s.
type ImageBundles struct {
	Images []ImageBundle `json:"images"`
}

// JSONShowPictures The JSON structure with offset and the number of pictures requested.
type JSONShowPictures struct {
	Offset int `json:"offset"`
	N int `json:"n"`
}

// JSONLabeledResult The JSON structure containing the name of the picture and the label created by the labeler.
// The label is a string and its meaning can be designed accordingly.
type JSONLabeledResult struct {
	Name string `json:"name"`
	Val string `json:"val"`
}

// JSONLabeledResults A collection of JSONLabeledResult s.
type JSONLabeledResults struct {
	Results []JSONLabeledResult `json:"results"`
}

// Returns a message string.
func completeMsgIfNameExist(pred *bool, label *bool, path *string) string {
	msg := "The picture is in the database.\n"
	if pred != nil {
		if *pred == true {
			msg += "prediction true\n"
		}else {
			msg += "prediction false\n"
		}
	} else {
		msg += "prediction unavailable\n"
	}
	if label != nil {
		if *label == true {
			msg += "label true\n"
		}else {
			msg += "label false\n"
		}
	} else {
		msg += "label unavailable\n"
	}
	msg += *path
	return msg
}

// Returns a message string
func completeMsgTensorflowResponse(pred bool) string {
	switch pred {
	case true:
		return "Prediction: true"
	case false:
		return "Prediction: false"
	default:
		return "completeMsgTensorflowResponse: Unknown prediction."
	}
}

// Gets the complete path of an image given the image name and the prediction.
func getPathPredicted(imgName string, pred bool) (string, error) {
	switch pred {
	case true:
		return common.S3VehiclePredictionPrefix + imgName, nil
	case false:
		return common.S3NonVehiclePredictionPrefix + imgName, nil
	default:
		return "", errors.New("getPathPredicted: Unknown prediction")
	}
}

// Closes the opened file.
func closeOpenedFile(file multipart.File) {
	err := file.Close()
	common.CheckErr(err)
}

// Closes the gRPC connection.
func closeGRPCConnection(conn *grpc.ClientConn) {
	err := conn.Close()
	common.CheckErr(err)
}

// Generate text in a special form. Used by "ShowPictures".
func generateText(id int, name string, prediction *bool, label *bool) string {
	r := "Id: " + strconv.Itoa(id) + ", name: " + name + ", prediction: "
	if prediction == nil {
		r = r + "None, label: "
	}else {
		r = r + strconv.FormatBool(*prediction) + ", label: "
	}
	if label == nil {
		r = r + "None."
	}else {
		r = r + strconv.FormatBool(*label) + "."
	}
	return r
}
