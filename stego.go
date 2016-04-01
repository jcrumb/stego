package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ltime)

	if len(os.Args) < 2 {
		log.Fatalln("Please provide a method")
	}

	if os.Args[1] == "-recover" {
		log.Printf("Recovery Mode")

		if len(os.Args) < 3 {
			log.Fatalln("Please provide a filename")
		}

		fileName := os.Args[2]
		recoverFile(fileName)

	} else if os.Args[1] == "-encrypt" {
		log.Printf("Encryption Mode")

		if len(os.Args) < 4 {
			log.Fatalln("Please follow the proper syntax -encrypt fileName message")
		}

		fileName := os.Args[2]
		message := os.Args[3]
		encrypt(fileName, message)

	} else if os.Args[1] == "-decrypt" {
		log.Printf("Decryption Mode")

		if len(os.Args) < 3 {
			log.Fatalln("Please provide a filename")
		}

		fileName := os.Args[2]
		decrypt(fileName)

	} else {
		log.Fatalln("Stego Help: -recover filename, -encrypt filename message, -decrypt filename")
	}

}

func encrypt(fileName string, message string) {
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

	outfile, err := os.Create("temp.bmp")
	checkError(err)
	defer outfile.Close()

	bmp.Encode(outfile, newImage)
}

func decrypt(fileName string) {
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

func recoverFile(fileName string) {
	log.Printf("Reading file %s into memory... ", fileName)
	fileData, err := ioutil.ReadFile(fileName)
	checkError(err)
	log.Printf("Done (read %d bytes)\n", len(fileData))

	imageData, err := recoverImages(fileData)
	checkError(err)

	err = os.Mkdir("recovered", 0777)
	checkError(err)

	for i := 0; i < len(imageData); i++ {
		path := fmt.Sprintf("recovered/image_%d.jpg", i)
		err = ioutil.WriteFile(path, imageData[i], 0777)
		checkError(err)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
