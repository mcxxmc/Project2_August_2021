package projectGo

import (
	"fmt"
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

type ImageBundle struct {
	EncodedImage string `json:"image"`
	Text string `json:"text"`
}

type ImageBundles struct {
	Images []ImageBundle `json:"images"`
}


// CheckErr checks the error and prints it out if not nil
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// CheckPanic checks the error and raise panic if needed
func CheckPanic(err error) {
	if err != nil {
		panic(err)
	}
}
