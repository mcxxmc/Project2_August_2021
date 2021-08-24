package projectGo

import (
	"bufio"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)

var PicTargetDir = "./pictures/"


// ReadImage reads the image given its name.
func ReadImage(imgPath string) image.Image {
	f, err := os.Open(imgPath)
	checkPanic(err)
	img, _, err := image.Decode(f)
	checkPanic(err)
	return img
}

// SaveImage saves the image.Image object on disk using the given name.
func SaveImage(imgName string, img image.Image) {
	outFile, err := os.Create(PicTargetDir + imgName)
	checkPanic(err)

	defer func(outFile *os.File) {  // close the outFile in the end
		err := outFile.Close()
		checkPanic(err)
	}(outFile)

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	checkPanic(err)

	err = b.Flush()
	checkPanic(err)
}
