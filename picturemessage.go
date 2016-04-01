package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"

	"golang.org/x/image/bmp"
)

func encodeMessage(fileName string, message string) {
	infile, err := os.Open(fileName)
	checkError(err)

	defer infile.Close()

	src, err := jpeg.Decode(infile)
	checkError(err)

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})

	numPixels := w * h
	pixelSpreader := (numPixels - 1) / len(message)
	counter := 1

	//Set first pixel R value to length of message
	oldColor := src.At(bounds.Min.X, bounds.Min.Y)
	_, initialG, initialB, _ := oldColor.RGBA()
	initialColor := color.RGBA{uint8(len(message)), uint8(initialG), uint8(initialB), 0}
	newImage.Set(bounds.Min.X, bounds.Min.Y, initialColor)

	//Spread out the pixels depending on the length of the string
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			if y == bounds.Min.Y && x == bounds.Min.X {
				continue
			}

			oldColor := src.At(x, y)
			r, g, b, _ := oldColor.RGBA()

			newColor := color.RGBA{uint8(r), uint8(g), uint8(b), 0}

			if (counter%pixelSpreader == 0 || counter == 1) && x*y != 0 {
				newColor = color.RGBA{uint8(message[counter/pixelSpreader-1]), uint8(g), uint8(b), 0}
			}

			newImage.Set(x, y, newColor)
			counter++

		}
	}

	outfile, err := os.Create(fileName)
	checkError(err)
	defer outfile.Close()

	bmp.Encode(outfile, newImage)
}

func decodeMessage(fileName string) {
	infile, err := os.Open(fileName)
	checkError(err)
	defer infile.Close()

	src, err := bmp.Decode(infile)
	checkError(err)

	bounds := src.Bounds()
	initialColor := src.At(bounds.Min.X, bounds.Min.Y)

	initialR, _, _, _ := initialColor.RGBA()
	stringLength := uint8(initialR)

	numPixels := bounds.Max.X * bounds.Max.Y
	pixelSpreader := (numPixels - 1) / int(stringLength)

	counter := 1
	hiddenMessage := ""

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if y == bounds.Min.Y && x == bounds.Min.X {
				continue
			}

			oldColor := src.At(x, y)
			r, _, _, _ := oldColor.RGBA()

			if (counter%pixelSpreader == 0 || counter == 1) && x*y != 0 {
				hiddenMessage += string(uint8(r))
			}
			counter++

		}
	}

	log.Printf("%s", hiddenMessage)
}
