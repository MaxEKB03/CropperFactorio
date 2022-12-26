package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
)

const NEW_DIR_NAME = "cropped/"
const DIR_NAME = "icons"

func CropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
    type subImager interface {
        SubImage(r image.Rectangle) image.Image
    }

    // img is an Image interface. This checks if the underlying value has a
    // method called SubImage. If it does, then we can use SubImage to crop the
    // image.
    simg, ok := img.(subImager)
    if !ok {
        return nil, fmt.Errorf("image does not support cropping")
    }

    return simg.SubImage(crop), nil
}

// Get square image from rectangle
// example 120:64 => 64:64
func CropMyFile(fileName string) image.Image{
	f, _ := os.Open(fileName)
    defer f.Close()
	oldImage, _ := png.Decode(f)

	bounds := oldImage.Bounds()
	size := bounds.Dy()-1

	rect := image.Rect(0, 0, size, size)

	newImage, _ := CropImage(oldImage, rect)

	return newImage
}

// Convert: image.Image => []byte
func GetBytesFromImage(newImage image.Image) []byte {
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, newImage)
	data := buf.Bytes()
	return data
}

// Recursion function for get all files
func ReadInDir(dirName string) {
	var dir, _ = os.ReadDir(dirName)
	os.Mkdir(NEW_DIR_NAME+dirName, fs.ModeDir)

	for _, file := range dir {
		fullFileName := dirName+"/"+file.Name()

		if file.IsDir() {
			ReadInDir(fullFileName)
		} else {
			newFileName := NEW_DIR_NAME+fullFileName

			newImage := CropMyFile(fullFileName)
			out, _ := os.Create(newFileName)
			png.Encode(out, newImage)
			out.Close()	
		}
	}
}

func main() {
	os.Mkdir(NEW_DIR_NAME, fs.ModeDir)
	ReadInDir(DIR_NAME)


	// catFile, _ := os.Open("cat.png")
    // defer catFile.Close()
 
    // cat, _ := png.Decode(catFile)

	// out, _ := os.Create("cat2.png")
	// png.Encode(out, cat)
	// out.Close()
}