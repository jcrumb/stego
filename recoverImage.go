package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func recoverImagesFromFile(fileName string) {
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

// recoverImage attempts to find JPEG files embedded in a byte slice, by looking for
// the FFD8 starting marker and reading up to an FFD9 ending marking if it exists
// If no image is found or the ending marker is missing, an error will be returned
// and the resultant byte slice will be nil
func recoverImages(fileData []byte) ([][]byte, error) {
	var images [][]byte
	var imageData []byte // slice for our eventual recovered image data

	startOffset := 0
	writing := false
	nestedSOI := 0

	for i := 0; i < len(fileData); i++ {
		if writing {
			imageData = append(imageData, fileData[i])
		}

		if fileData[i] == byte(JFIF_STARTMARK) {
			switch fileData[i+1] {
			case byte(JFIF_SOI):
				if writing {
					// already writing, so we must be seeing the start of a nested jpg
					nestedSOI++
				} else {
					writing = true
					imageData = append(imageData, fileData[i])
					log.Printf("Found JPG starting at offset %d", i)
					startOffset = i
				}
			case byte(JFIF_EOI):
				if nestedSOI == 0 {
					// End of a non nested image
					imageData = append(imageData, fileData[i+1])

					log.Printf("Found end of JPG at offset %d", i)
					log.Printf("Total size: %d bytes", i-startOffset)

					images = append(images, imageData)

					writing = false
					i++
					imageData = make([]byte, 0)

				} else {
					nestedSOI--
				}
			}
		}

	}
	log.Printf("Found %d total images", len(images))
	return images, nil
}
