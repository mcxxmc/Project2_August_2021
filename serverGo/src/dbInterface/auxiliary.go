package dbInterface

import (
	"strconv"
)

// Record A complete record (including id, name, path, prediction, label) in the database.
type Record struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Prediction *bool `json:"prediction"`
	Label *bool `json:"label"`
}

// Records A collection of Record s.
type Records struct {
	Recs []Record `json:"records"`
}

// PathAndDesc The JSON structure containing the path of an image and a customized description of that image.
type PathAndDesc struct {
	Path string
	Text string
}

// UnlabeledRecord The JSON structure containing the information (name and path) of an unlabeled image.
type UnlabeledRecord struct {
	Name string
	Path string
}

type UnpredictedUnlabeledRecord struct {
	Name string
	Path string
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
