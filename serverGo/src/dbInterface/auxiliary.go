package dbInterface

import (
	"fmt"
	"strconv"
)

type Record struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Prediction *bool `json:"prediction"`
	Label *bool `json:"label"`
}

type Records struct {
	Recs []Record `json:"records"`
}

type JSONShowPictures struct {
	Offset int `json:"offset"`
	N int `json:"n"`
}

type JSONLabeledResult struct {
	Name string `json:"name"`
	Val string `json:"val"`
}

type JSONLabeledResults struct {
	Results []JSONLabeledResult `json:"results"`
}

type PathAndDesc struct {
	Path string
	Text string
}

type UnlabeledRecord struct {
	Name string
	Path string
}

type ImageBundle struct {
	EncodedImage string `json:"image"`
	Text string `json:"text"`
}

type ImageBundles struct {
	Images []ImageBundle `json:"images"`
}

// checkErr checks the error and prints it out if not nil
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// generate text for ShowPictures
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
