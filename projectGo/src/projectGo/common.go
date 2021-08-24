package projectGo

import (
	"fmt"
	"image"
)

type BufferedImage struct {
	ImgName string
	Img     image.Image
}

// CheckErr checks the error and prints it out if not nil
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// checks if the error is not nil, and raise a panic if needed
func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func generateImgPath(imgName string) string {
	return PicTargetDir + imgName
}

func mapInt2Bool(n int) bool {
	if n == 0 {
		return true
	}
	return false
}

func mapBool2Int(b bool) int {
	if b == true {
		return 0
	}
	return 1
}
