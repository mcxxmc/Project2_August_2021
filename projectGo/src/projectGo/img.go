package projectGo

import (
	"bufio"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)


// LoadImage loads the image given its path.
func LoadImage(imgPath string) image.Image {
	f, err := os.Open(imgPath)
	CheckPanic(err)
	img, _, err := image.Decode(f)
	CheckPanic(err)
	return img
}

// SaveImage saves the image.Image object on disk using the given name.
// Warning: the first parameter should be a valid path! (a complete path is recommended)
func SaveImage(imgPath string, img image.Image) {
	outFile, err := os.Create(imgPath)
	CheckPanic(err)

	defer func(outFile *os.File) {  // close the outFile in the end
		err := outFile.Close()
		CheckPanic(err)
	}(outFile)

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	CheckPanic(err)

	err = b.Flush()
	CheckPanic(err)
}
