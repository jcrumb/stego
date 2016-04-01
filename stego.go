package main

import (
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
		recoverImagesFromFile(fileName)

	} else if os.Args[1] == "-encrypt" {
		log.Printf("Encryption Mode")

		if len(os.Args) < 4 {
			log.Fatalln("Please follow the proper syntax -encrypt fileName message")
		}

		fileName := os.Args[2]
		message := os.Args[3]
		encodeMessage(fileName, message)

	} else if os.Args[1] == "-decrypt" {
		log.Printf("Decryption Mode")

		if len(os.Args) < 3 {
			log.Fatalln("Please provide a filename")
		}

		fileName := os.Args[2]
		decodeMessage(fileName)

	} else {
		log.Fatalln("Stego Help: -recover filename, -encrypt filename message, -decrypt filename")
	}

}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
