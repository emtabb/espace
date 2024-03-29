package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func storeImage(imgByte []byte, path string) image.Image {
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := os.Create(path)
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 1

	err = jpeg.Encode(out, img, &opts)

	if err != nil {
		log.Println(err)
	}
	return img
}

func generateImage() {
	// Create a blank image 10 pixels wide by 4 pixels tall
	myImage := image.NewRGBA(image.Rect(0, 0, 10, 4))

	// You can access the pixels through myImage.Pix[i]
	// One pixel takes up four bytes/uint8. One for each: RGBA
	// So the first pixel is controlled by the first 4 elements
	// Values for color are 0 black - 255 full color
	// Alpha value is 0 transparent - 255 opaque
	myImage.Pix[0] = 255 // 1st pixel red
	myImage.Pix[1] = 0 // 1st pixel green
	myImage.Pix[2] = 0 // 1st pixel blue
	myImage.Pix[3] = 255 // 1st pixel alpha

	// myImage.Pix contains all the pixels
	// in a one-dimensional slice
	fmt.Println(myImage.Pix)

	// Stride is how many bytes take up 1 row of the image
	// Since 4 bytes are used for each pixel, the stride is
	// equal to 4 times the width of the image
	// Since all the pixels are stored in a 1D slice,
	// we need this to calculate where pixels are on different rows.
	fmt.Println(myImage.Stride) // 40 for an image 10 pixels wide
}

func readImage(path string) (image.Image, error) {
	existingImageFile, err := os.Open(path)
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	imageData, imageType, err := image.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}
	fmt.Println(imageData)
	fmt.Println(imageType)

	// We only need this because we already read from the file
	// We have to reset the file pointer back to beginning
	existingImageFile.Seek(0, 0)

	// Alternatively, since we know it is a png already
	// we can call png.Decode() directly
	return jpeg.Decode(existingImageFile)
}
