package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ltime)

	if len(os.Args) != 2 {
		log.Fatalln("Please provide a filename")
	}
	fileName := os.Args[1]

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
